package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/ESSantana/streaming-test/internal/services/interfaces"
)

type VideoData struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type VideoUploader struct {
	videoService interfaces.VideoService
}

func NewVideoUploader(videoService interfaces.VideoService) *VideoUploader {
	return &VideoUploader{
		videoService: videoService,
	}
}

func (v *VideoUploader) CreateS3PresignedPutURL(w http.ResponseWriter, r *http.Request) {
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

	url, err := v.videoService.CreateS3PresignedPutURL(r.Context(), os.Getenv("VIDEO_BUCKET"), videoData.Filename, videoData.ContentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url))
}
