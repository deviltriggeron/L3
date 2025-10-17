package router

import (
	"commentTree/internal/handler"
	"commentTree/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(h *handler.CommetsHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.Logger)

	r.HandleFunc("/", h.Index)
	r.HandleFunc("/allComments", h.GetAllParentComments).Methods("GET")
	r.HandleFunc("/comments", h.NewComments).Methods("POST")
	r.HandleFunc("/comments", h.GetComments).Methods("GET")
	r.HandleFunc("/comments/{id}", h.DeleteComments).Methods("DELETE")

	return r
}
