package dto

import "github.com/ESSantana/streaming-test/internal/domain/models"

type VideoUploadRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type VideoUploadResponse struct {
	UploadURL string `json:"upload_url"`
}

type ListVideosResponse struct {
	Videos []models.Video `json:"videos"`
}

type VideoDistributionResponse struct {
	VideoURL string `json:"video_url"`
}
