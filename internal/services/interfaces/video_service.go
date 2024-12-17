package interfaces

import "context"

type VideoService interface {
	UploadVideo(ctx context.Context, filename string, extension string, data []byte) (err error)
}
