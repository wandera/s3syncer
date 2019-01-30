package sync

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// main
func RunSync(folderToWatch, s3Region, s3Path string) error {
	s3Path = strings.Replace(s3Path, "s3://", "", -1)
	s3Bucket := strings.Split(s3Path, "/")[0]
	s3Folder := strings.Join(strings.Split(s3Path, "/")[1:], "/")
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)
				// watch for errors
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Info("modified file:", event.Name, "\nUploading to S3")
					uploadFile(event.Name, s3Region, s3Bucket, s3Folder)
				case event.Op&fsnotify.Create == fsnotify.Create:
					log.Info("New file:", event.Name, "\nUploading to S3")
					uploadFile(event.Name, s3Region, s3Bucket, s3Folder)
				default:
					log.Println("No changes")
				}
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// Add the folder to watch to the Watcher
	if err := watcher.Add(folderToWatch); err != nil {
		return err
	}

	<-done
	return nil
}

func uploadFile(fileDir, s3Region, s3Bucket, s3Folder string) {
	// Create a single AWS session (we can re use this if we're uploading many files)
	s, err := session.NewSession(&aws.Config{Region: aws.String(s3Region)})
	if err != nil {
		log.Error(err)
	}

	// Upload
	err = addFileToS3(s, fileDir, s3Bucket, s3Folder)
	if err != nil {
		log.Error(err)
	}
}

// addFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func addFileToS3(s *session.Session, fileDir, s3Bucket, s3Folder string) error {

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	_, _ = file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(s3Bucket),
		Key:                aws.String(s3Folder + filepath.Base(fileDir)),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	return err
}
