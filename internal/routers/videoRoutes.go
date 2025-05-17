package routers

import (
	"github.com/ESSantana/streaming-test/internal/controllers"
	iservice "github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/go-chi/chi/v5"
)

func configureVideoRoutes(router *chi.Mux, serviceManager iservice.ServiceManager) {
	videoController := controllers.NewVideoController(serviceManager)

	router.Post("/upload", videoController.CreateS3PresignedPutURL)
	router.Post("/video-processor", videoController.ProcessVideo)
	router.Get("/videos", videoController.ListAvailableVideos)
}
