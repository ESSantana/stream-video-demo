package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ESSantana/streaming-test/internal/domain"
	iservice "github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/ESSantana/streaming-test/pkg/dto"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/rs/zerolog/log"
)

type VideoController struct {
	serviceManager iservice.ServiceManager
}

func NewVideoController(serviceManager iservice.ServiceManager) *VideoController {
	return &VideoController{
		serviceManager: serviceManager,
	}
}

func (ctrl *VideoController) UploadVideo(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var videoUploadRequest dto.VideoUploadRequest
	err = json.Unmarshal(body, &videoUploadRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	videoService := ctrl.serviceManager.NewVideoService()

	uploadURL, err := videoService.UploadRawVideo(r.Context(), videoUploadRequest.Filename, videoUploadRequest.ContentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := dto.VideoUploadResponse{
		UploadURL: url.QueryEscape(uploadURL),
	}
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (ctrl *VideoController) ListAvailableVideos(w http.ResponseWriter, r *http.Request) {
	videoService := ctrl.serviceManager.NewVideoService()

	availableVideos, err := videoService.ListAvailableVideos(r.Context(), os.Getenv("VIDEO_BUCKET"))
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

func (ctrl *VideoController) ProcessRawVideo(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if strings.Contains(string(data), "ConfirmSubscription") {
		ctrl.confirmSNSSubscription(data)
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

	if len(s3Events.Records) < 1 { 
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Any video to process"))
	}

	videoService := ctrl.serviceManager.NewVideoService()
	videoKey := s3Events.Records[0].S3.Object.Key
	go func() {
		err := videoService.ProcessVideoWithOptions(r.Context(), videoKey, domain.DefaultVideoOptions())
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processing video"))
}

func (ctrl *VideoController) confirmSNSSubscription(data []byte) {
	var subscriptionInput sns.ConfirmSubscriptionInput
	err := json.Unmarshal(data, &subscriptionInput)
	if err != nil {
		//Good to have: notify the user that the subscription could not be confirmed
		log.Error().Msg(err.Error())
		return
	}

	session, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	snsClient := sns.New(session, aws.NewConfig().WithRegion("sa-east-1"))

	_, err = snsClient.ConfirmSubscription(&subscriptionInput)
	if err != nil {
		//Good to have: notify the user that the subscription could not be confirmed
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msg("SNS subscription confirmed")
}
