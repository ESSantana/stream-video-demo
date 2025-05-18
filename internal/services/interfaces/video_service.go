package interfaces

import (
	"context"

	"github.com/ESSantana/streaming-test/internal/domain"
)

type VideoService interface {
	UploadRawVideo(ctx context.Context, filename, contentType string) (presignedURL string, err error)
	ProcessVideoWithOptions(ctx context.Context, videoKey string, options domain.VideoOptions) (err error)
	ListAvailableVideos(ctx context.Context, bucket string) (availableVideos []string, err error) 
}
