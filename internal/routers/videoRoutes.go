package routers

import (
	"github.com/ESSantana/streaming-test/internal/controllers"
	iservice "github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/go-chi/chi/v5"
)

func configureVideoRoutes(router *chi.Mux, serviceManager iservice.ServiceManager) {
	videoController := controllers.NewVideoController(serviceManager)

	router.Post("/video", videoController.UploadVideo)
	router.Get("/video", videoController.ListAvailableVideos)
	router.Post("/video-processor", videoController.ProcessRawVideo)
}
