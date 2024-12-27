package handlers

import (
	"net/http"
)

type VideoProcessor struct {
}

func NewVideoProcessor() *VideoProcessor {
	return &VideoProcessor{}
}

func (v *VideoProcessor) Process(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, world!"))
}
