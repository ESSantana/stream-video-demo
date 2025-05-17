package storage

import (
	"fmt"
	"io"
	"os"
	"time"

	istorage "github.com/ESSantana/streaming-test/internal/storage/interfaces"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type storageManager struct {
	client *s3.S3
	bucket string
}

func NewStorageManager(client *s3.S3) istorage.StorageManager {
	return &storageManager{
		client: client,
		bucket: os.Getenv("VIDEO_BUCKET"),
	}
}

func (s *storageManager) UploadRawVideo(filename, contentType string) (uploadURL string, err error) {
	objectKey := "raw/" + filename

	req, _ := s.client.PutObjectRequest(
		&s3.PutObjectInput{
			Bucket:      aws.String(s.bucket),
			Key:         aws.String(objectKey),
			ContentType: aws.String(contentType),
		},
	)

	uploadURL, err = req.Presign(time.Minute * 15)
	return uploadURL, err
}

func (s *storageManager) UploadProcessedVideo(basePath, videoID string, processedFiles []os.DirEntry) (err error) {
	for _, processedFile := range processedFiles {
		tempFilePath := basePath + processedFile.Name()

		data, err := os.OpenFile(tempFilePath, os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		objectKey := fmt.Sprintf("processed/%s/%s", videoID, processedFile.Name())
		_, err = s.client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(objectKey),
			Body:   data,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *storageManager) RetrieveRawVideo(objectKey string) (objectData []byte, err error) {
	videoObject, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	defer videoObject.Body.Close()

	objectData, err = io.ReadAll(videoObject.Body)
	return objectData, err
}
