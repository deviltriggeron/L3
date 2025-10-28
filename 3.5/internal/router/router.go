package router

import (
	"github.com/gorilla/mux"

	"eventbooker/internal/handler"
	"eventbooker/internal/middleware"
)

func NewEventBookRouter(h *handler.EventBookHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/", h.Index)
	r.HandleFunc("/events", h.CreateEvent).Methods("POST")
	r.HandleFunc("/users", h.CreateUser).Methods("POST")
	r.HandleFunc("/events/{id}/book", h.Book).Methods("POST")
	r.HandleFunc("/events/{id}/confirm", h.Payment).Methods("POST")
	r.HandleFunc("/events/{id}", h.GetEvent).Methods("GET")

	return r
}
