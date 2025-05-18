package routers

import (
	"encoding/json"
	"net/http"
	"time"

	iservice "github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/go-chi/chi/v5"
)

func ConfigureRouter(router *chi.Mux, serviceManager iservice.ServiceManager) {
	router.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(map[string]string{"health": "ok", "time": time.Now().Format(time.RFC3339)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	configureVideoRoutes(router, serviceManager)
}
