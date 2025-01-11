package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/ESSantana/streaming-test/pkg/dto"
	"github.com/go-chi/chi/v5"
)

type VideoController struct {
	videoService  interfaces.VideoService
	defaultClient *http.Client
}

func NewVideoController(videoService interfaces.VideoService) *VideoController {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}
	return &VideoController{
		videoService:  videoService,
		defaultClient: &httpClient,
	}
}

func (v *VideoController) CreateS3PresignedPutURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var videoData dto.VideoUploadRequest
	err = json.Unmarshal(body, &videoData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	presignedURL, err := v.videoService.CreateS3PresignedPutURL(r.Context(), os.Getenv("VIDEO_BUCKET"), videoData.Filename, videoData.ContentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := dto.VideoUploadResponse{
		PresignedURL: url.QueryEscape(presignedURL),
	}
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (v *VideoController) ListAvailableVideos(w http.ResponseWriter, r *http.Request) {
	availableVideos, err := v.videoService.ListAvailableVideos(r.Context(), os.Getenv("VIDEO_BUCKET"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(availableVideos) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data := dto.ListVideosResponse{
		Videos: availableVideos,
	}
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (v *VideoController) GetVideoDistribution(w http.ResponseWriter, r *http.Request) {
	videoName := chi.URLParam(r, "video")

	videoDistributionURL := v.mountAndValidateDistributionURL(videoName)
	if videoDistributionURL == "" {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	data := dto.VideoDistributionResponse{
		VideoURL: videoDistributionURL,
	}
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (v *VideoController) mountAndValidateDistributionURL(videoName string) string {
	videoDistributionURL := "https://" + os.Getenv("CLOUDFRONT_DIST") + "/" + "processed/" + videoName + "/index.m3u8"

	res, err := v.defaultClient.Head(videoDistributionURL)
	if err != nil {
		return ""
	}

	if res.StatusCode != http.StatusOK {
		return ""
	}

	return videoDistributionURL
}
