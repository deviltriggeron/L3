package router

import (
	"shortener/internal/handler"

	middleware "shortener/internal/logger"

	"github.com/gorilla/mux"
)

/*
– POST /shorten — создание новой сокращённой ссылки;
– GET /s/{short_url} — переход по короткой ссылке;
– GET /analytics/{short_url} — получение аналитики (число переходов, User-Agent, время переходов).
*/

func NewRouter(h *handler.ShortenerHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.Logger)

	r.HandleFunc("/shorten", h.NewShorten).Methods("POST")
	r.HandleFunc("/s/{short_url}", h.Redirect).Methods("GET")
	r.HandleFunc("/analytics/{short_url}", h.Analytics).Methods("GET")

	return r
}
