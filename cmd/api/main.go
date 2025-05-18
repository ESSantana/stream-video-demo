package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ESSantana/streaming-test/internal/repositories"
	irepository "github.com/ESSantana/streaming-test/internal/repositories/interfaces"
	"github.com/ESSantana/streaming-test/internal/routers"
	"github.com/ESSantana/streaming-test/internal/services"
	iservices "github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/ESSantana/streaming-test/internal/storage"
	istorage "github.com/ESSantana/streaming-test/internal/storage/interfaces"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)

var (
	router         *chi.Mux
	repositoryManager irepository.RepositoryManager
	serviceManager iservices.ServiceManager
	storageManager istorage.StorageManager
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
	routers.ConfigureRouter(router, serviceManager)
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
	repositoryManager, err = repositories.NewRepositoryManager(context.Background())
	if err != nil { 
		panic(err)
	}
	storageManager = storage.NewStorageManager(s3Client)
	serviceManager = services.NewServiceManager(storageManager, repositoryManager)
}
 