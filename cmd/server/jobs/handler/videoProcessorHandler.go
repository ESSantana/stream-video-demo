package handler

import (
	"fmt"
	"io"
	"net/http"
	// "os"

	"github.com/ESSantana/streaming-test/internal/services/interfaces"
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

	fmt.Println(string(data))

	// err := h.videoService.ProcessVideoWithOptions(r.Context(), os.Getenv("VIDEO_BUCKET"), "videoKey", nil)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}
