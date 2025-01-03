package main

import (
	"context"

	"github.com/ESSantana/streaming-test/cmd/serverless/api/controllers"
	"github.com/ESSantana/streaming-test/internal/services"
	"github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	chiLambda *chiadapter.ChiLambda
	videoService interfaces.VideoService
)

func init() {
	loadDependencies()

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	videoController := controllers.NewVideoUploader(videoService)
	router.Post("/upload", videoController.CreateS3PresignedPutURL)

	chiLambda = chiadapter.New(router)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func loadDependencies() {
	session, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s3Client := s3.New(session, aws.NewConfig().WithRegion("sa-east-1"))

	videoService = services.NewVideoService(s3Client)
}
