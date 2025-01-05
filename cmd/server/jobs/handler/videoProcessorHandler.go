package handler

import (
	"io"
	"net/http"
	// "os"

	"github.com/ESSantana/streaming-test/internal/services/interfaces"
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

	// go func() {
	// 	err = h.videoService.ProcessVideoWithOptions(r.Context(), os.Getenv("VIDEO_BUCKET"), "raw/epic_sax_guy.mp4", nil)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	if err != nil {
	// 		log.Error().Msg(err.Error())
	// 	}
	// }()

	log.Info().Msg(string(data))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processing video"))
}
