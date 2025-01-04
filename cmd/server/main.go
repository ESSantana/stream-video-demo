package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ESSantana/streaming-test/cmd/server/jobs/handler"
	"github.com/ESSantana/streaming-test/internal/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

var (
	router                *chi.Mux
	videoProcessorHandler *handler.VideoProcessorHandler
)

func main() {
	loadDependencies()
	setupRoute()
	defer startServer(router)
	fmt.Printf("Server listening on port :%s\n", os.Getenv("SERVER_PORT"))
}

func setupRoute() {
	router = chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	router.Post("/video-processor", videoProcessorHandler.ProcessVideo)
}

func startServer(router *chi.Mux) {
	port := ":" + os.Getenv("SERVER_PORT")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	if err := http.Serve(listen, router); err != nil {
		panic(err)
	}
}

func loadDependencies() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	session, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	s3Client := s3.New(session, aws.NewConfig().WithRegion("sa-east-1"))
	videoService := services.NewVideoService(s3Client)

	videoProcessorHandler = handler.NewVideoProcessorHandler(videoService)
}
