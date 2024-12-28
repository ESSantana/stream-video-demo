package main

import (
	"context"
	"net/http"

	"github.com/ESSantana/streaming-test/cmd/serverless/api/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var chiLambda *chiadapter.ChiLambda

func init() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	session, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s3Client := s3.New(session, aws.NewConfig().WithRegion("sa-east-1"))

	videoUploader := handlers.NewVideoUploader(s3Client)

	router.Get("/ping", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("pong"))
	})

	router.Post("/upload", videoUploader.Process)

	chiLambda = chiadapter.New(router)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}
