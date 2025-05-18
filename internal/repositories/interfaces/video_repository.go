package interfaces

import (
	"context"

	"github.com/ESSantana/streaming-test/internal/domain/models"
)

type VideoRepository interface {
	Save(ctx context.Context, video models.Video) (err error)
	ListAvailableVideos(ctx context.Context) (videos []models.Video, err error)
}
