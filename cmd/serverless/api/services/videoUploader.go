package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	// "github.com/google/uuid"
)

type VideoUploader struct {
	s3Client *s3.S3
}

func NewVideoUploader(s3Client *s3.S3) *VideoUploader {
	return &VideoUploader{
		s3Client: s3Client,
	}
}

func (v *VideoUploader) Process(w http.ResponseWriter, r *http.Request) {
	// fileID := uuid.New().String()

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	videoFile, videoHeader, err := r.FormFile("upload_video")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer videoFile.Close()

	fmt.Printf("File Name: %s, Size: %v", videoHeader.Filename, videoHeader.Size)

	data, err := io.ReadAll(videoFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(data) < 1 {
		http.Error(w, "Error at upload video", http.StatusBadRequest)
		return
	}

	req, _ := v.s3Client.PutObjectRequest(
		&s3.PutObjectInput{
			Bucket: aws.String("streaming-test-essantana"),
			Key:    aws.String("videozinho"),
			Body:   strings.NewReader("EXPECTED CONTENTS"),
		},
	)

	url, err := req.Presign(time.Minute * 15)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url))

	// err = os.MkdirAll(tempDir, os.ModePerm)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// tempFilePath := fmt.Sprintf("%s/%s.mp4", tempDir, fileID)
	// err = os.WriteFile(tempFilePath, data, os.ModeAppend)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}
