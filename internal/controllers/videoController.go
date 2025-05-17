package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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
	defaultClient  *http.Client
}

func NewVideoController(serviceManager iservice.ServiceManager) *VideoController {
	//TODO: receive as dependency
	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	return &VideoController{
		serviceManager: serviceManager,
		defaultClient:  &httpClient,
	}
}

// TODO: Update name to be tech agnostic
func (ctrl *VideoController) CreateS3PresignedPutURL(w http.ResponseWriter, r *http.Request) {
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

	//TODO: Remove bucket name dependency
	uploadURL, err := videoService.CreateS3PresignedPutURL(r.Context(), videoUploadRequest.Filename, videoUploadRequest.ContentType)
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
	//TODO: Remove bucket name dependency
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

// TODO: Update name to be more descritive
func (ctrl *VideoController) ProcessVideo(w http.ResponseWriter, r *http.Request) {
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

	for _, record := range s3Events.Records {
		go ctrl.createProcessingRoutine(r.Context(), record.S3.Object.Key)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processing video"))
}

// TODO: Move it to another place
func (ctrl *VideoController) createProcessingRoutine(ctx context.Context, videoKey string) {
	videoService := ctrl.serviceManager.NewVideoService()
	//TODO: Remove bucket name dependency
	err := videoService.ProcessVideoWithOptions(ctx, videoKey, domain.DefaultVideoOptions())
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

// TODO: Move it to another place
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
