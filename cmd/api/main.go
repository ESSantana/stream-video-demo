package main

import (
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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)

var (
	router            *chi.Mux
	repositoryManager irepository.RepositoryManager
	serviceManager    iservices.ServiceManager
	storageManager    istorage.StorageManager
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

	var err error
	repositoryManager, err = repositories.NewRepositoryManager()
	if err != nil {
		panic(err)
	}
	storageManager, err = storage.NewStorageManager()
	if err != nil {
		panic(err)
	}
	serviceManager = services.NewServiceManager(storageManager, repositoryManager)
}
