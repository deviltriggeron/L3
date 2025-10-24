package router

import (
	"github.com/gorilla/mux"

	"imageprocessor/internal/handler"
	"imageprocessor/internal/middleware"
)

func NewRouter(h *handler.ImageProcHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/", h.Index)
	r.HandleFunc("/upload", h.Upload).Methods("POST")
	r.HandleFunc("/image/{id}", h.GetImage).Methods("GET")
	r.HandleFunc("/image/{id}", h.DeleteImage).Methods("DELETE")

	return r
}
