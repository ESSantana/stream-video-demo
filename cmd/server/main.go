package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ESSantana/streaming-test/cmd/server/api/controllers"
	"github.com/ESSantana/streaming-test/cmd/server/jobs/handler"
	"github.com/ESSantana/streaming-test/internal/services"
	iservices "github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)

var (
	router       *chi.Mux
	videoService iservices.VideoService
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
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(map[string]string{"health": "ok", "time": time.Now().Format(time.RFC3339)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	// API endpoints
	videoController := controllers.NewVideoController(videoService)
	router.Post("/upload", videoController.CreateS3PresignedPutURL)
	router.Get("/videos", videoController.ListAvailableVideos)
	router.Get("/videos/{video}", videoController.GetVideoDistribution)

	// Jobs endpoints
	videoProcessorHandler := handler.NewVideoProcessorHandler(videoService)
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
	videoService = services.NewVideoService(s3Client)
}
