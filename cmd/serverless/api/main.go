package api

import (
	"context"
	"log"
	"net/http"

	"github.com/ESSantana/streaming-test/cmd/serverless/api/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var chiLambda *chiadapter.ChiLambda

func init() {
	log.Printf("Chi cold start")
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	videoUploader := handlers.NewVideoUploader()

	router.Get("/ping", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("pong"))
	})

	router.Post("/upload", videoUploader.Process)

	chiLambda = chiadapter.New(router)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
