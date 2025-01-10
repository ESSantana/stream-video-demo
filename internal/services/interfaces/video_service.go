package interfaces

import (
	"context"

	"github.com/ESSantana/streaming-test/internal/domain"
)

type VideoService interface {
	CreateS3PresignedPutURL(ctx context.Context, bucket, filename, contentType string) (presignedURL string, err error)
	ProcessVideoWithOptions(ctx context.Context, bucket, videoKey string, options domain.VideoOptions) (err error)
	ListAvailableContent(ctx context.Context, bucket string) (availableContent []string, err error) 
}
