package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type VideoData struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type VideoUploader struct {
	s3Client *s3.S3
}

func NewVideoUploader(s3Client *s3.S3) *VideoUploader {
	return &VideoUploader{
		s3Client: s3Client,
	}
}

func (v *VideoUploader) Process(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var videoData VideoData
	err = json.Unmarshal(body, &videoData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, _ := v.s3Client.PutObjectRequest(
		&s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("VIDEO_BUCKET")),
			Key:    aws.String("raw/" + videoData.Filename),
			ContentType: aws.String(videoData.ContentType),
		},
	)

	url, err := req.Presign(time.Minute * 15)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url))
}
