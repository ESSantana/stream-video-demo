package interfaces

import "context"

type VideoService interface {
	CreateS3PresignedPutURL(ctx context.Context, bucket, filename, contentType string) (presignedURL string, err error)
	ProcessVideoWithOptions(ctx context.Context, bucket, videoKey string, options interface{}) (err error)
}
