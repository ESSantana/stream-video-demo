package interfaces

import (
	"context"

	"github.com/ESSantana/streaming-test/internal/domain"
	"github.com/ESSantana/streaming-test/internal/domain/models"
)

type VideoService interface {
	UploadRawVideo(ctx context.Context, filename, contentType string) (presignedURL string, err error)
	ProcessVideoWithOptions(ctx context.Context, videoKey string, options domain.VideoOptions) (err error)
	ListAvailableVideos(ctx context.Context) (availableVideos []models.Video, err error)
}
