package handler

import (
	"context"
	"os"

	"github.com/ESSantana/streaming-test/internal/services/interfaces"
)

type VideoProcessorHandler struct {
	videoService interfaces.VideoService
}

func NewVideoProcessorHandler(videoService interfaces.VideoService) *VideoProcessorHandler {
	return &VideoProcessorHandler{videoService: videoService}
}

func (h *VideoProcessorHandler) ProcessVideo(ctx context.Context, videoKey string) error {
	return h.videoService.ProcessVideoWithOptions(ctx, os.Getenv("VIDEO_BUCKET"), videoKey, nil)

}
