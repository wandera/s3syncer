package sync

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Syncer interface {
	Sync() error
}

type syncer struct {
	folderToWatch string
	s3Bucket      string
	s3Folder      string
	stop          chan struct{}
	watcher       *fsnotify.Watcher
	s3session     *session.Session
	uploader      *s3manager.Uploader
}

func NewSyncer(folderToWatch, s3Region, s3Path string, stop chan struct{}) (Syncer, error) {
	s3Path =  strings.Replace(s3Path, "s3://", "", -1)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	s3session, err := session.NewSession(&aws.Config{Region: aws.String(s3Region)})
	if err != nil {
		return nil, err
	}
	uploader := s3manager.NewUploader(s3session)
	return &syncer{
		folderToWatch: folderToWatch,
		s3Bucket:      strings.Split(s3Path, "/")[0],
		s3Folder:      strings.Join(strings.Split(s3Path, "/")[1:], "/"),
		stop:          stop,
		watcher:       watcher,
		s3session:     s3session,
		uploader:      uploader,
	}, nil
}

func (s *syncer) Sync() error {
	go func() {
		for {
			select {
			// watch for events
			case event := <-s.watcher.Events:
				log.Debug(event)
				// watch for errors
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Infof("uploading modified file '%s' to S3", event.Name)
					go s.uploadFile(event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					log.Infof("uploading new file '%s' to S3", event.Name)
					go s.uploadFile(event.Name)
				default:
					log.Debug("no changes")
				}
			case err := <-s.watcher.Errors:
				log.Error(err)
			case <-s.stop:
				s.watcher.Close()
			}
		}
	}()

	// Add the folder to watch to the Watcher
	if err := s.watcher.Add(s.folderToWatch); err != nil {
		return err
	}

	return nil
}

func (s *syncer) uploadFile(fileDir string) {
	// Upload
	err := s.addFileToS3(s.s3session, fileDir)
	if err != nil {
		log.Error(err)
	}
}

// addFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (s *syncer) addFileToS3(sess *session.Session, fileDir string) error {

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(s.s3Bucket),
		Key:                aws.String(s.s3Folder + filepath.Base(fileDir)),
		ACL:                aws.String(s3.ObjectCannedACLPrivate),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
		Body:               file,
	})
	return err
}
