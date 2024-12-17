package main

import (
	"context"
	"net/http"

	"github.com/ESSantana/streaming-test/internal/services"
	iservices "github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

var (
	videoService iservices.VideoService
)

func init() {
	videoService = services.NewVideoService()
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	data := []byte(req.Body)
	if data == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error at upload video",
		}, nil
	}

	fileID := uuid.New().String()

	err = videoService.UploadVideo(ctx, fileID, ".mp4", data)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "successfully uploaded video",
	}, nil
}
