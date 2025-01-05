package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

type VideoProcessorHandler struct {
	videoService interfaces.VideoService
}

func NewVideoProcessorHandler(videoService interfaces.VideoService) *VideoProcessorHandler {
	return &VideoProcessorHandler{videoService: videoService}
}

func (h *VideoProcessorHandler) ProcessVideo(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if os.Getenv("SUBSCRIBE_TOPIC_MODE") == "true" {
		log.Info().Msg(string(data))
		return
	}

	var snsMessage events.SNSEntity
	err = json.Unmarshal(data, &snsMessage)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	var s3Events events.S3Event
	err = json.Unmarshal([]byte(snsMessage.Message), &s3Events)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	for _, record := range s3Events.Records {
		go h.createProcessingRoutine(r.Context(), record.S3.Object.Key)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processing video"))
}

func (h *VideoProcessorHandler) createProcessingRoutine(ctx context.Context, videoKey string) {
	err := h.videoService.ProcessVideoWithOptions(ctx, os.Getenv("VIDEO_BUCKET"), videoKey, nil)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
