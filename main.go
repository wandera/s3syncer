package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var s3Region, s3Bucket, s3Folder string

// main
func main() {
	folderToWatch := os.Args[1]
	s3Path := os.Args[2]
	s3Bucket = strings.Split(s3Path, "/")[0]
	s3Folder = strings.Join(strings.Split(s3Path, "/")[1:], "/")
	s3Region = os.Args[3]
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
					log.Println("modified file:", event.Name, "\nUploading to S3")
					uploadFile(event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					log.Println("New file:", event.Name, "\nUploading to S3")
					uploadFile(event.Name)
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
		fmt.Println("ERROR", err)
	}

	<-done
}

func uploadFile(fileDir string) {

	// Create a single AWS session (we can re use this if we're uploading many files)
	s, err := session.NewSession(&aws.Config{Region: aws.String(s3Region)})
	if err != nil {
		log.Println("ERROR", err)
	}

	// Upload
	err = AddFileToS3(s, fileDir)
	if err != nil {
		log.Println("ERROR", err)
	}
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, fileDir string) error {

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
