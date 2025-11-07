package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"warehouse-control/internal/handler"
	"warehouse-control/internal/interfaces"
	"warehouse-control/internal/middleware"
)

func NewRouter(handler handler.ControllerHandler, authHandler handler.AuthHandler, historyHandler handler.HistoryHandler, tp interfaces.TokenProvide, db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.MiddlewareLogging)

	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(func(next http.Handler) http.Handler {
		return middleware.JWTAuth(tp, next)
	})
	api.Use(middleware.WithDBUser(db))

	api.HandleFunc("/items", handler.AddItem).Methods("POST")
	api.HandleFunc("/items", handler.GetAllItem).Methods("GET")
	api.HandleFunc("/items/{id}", handler.GetItem).Methods("GET")
	api.HandleFunc("/items/{id}", handler.UpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", handler.DeleteItem).Methods("DELETE")

	api.HandleFunc("/history", historyHandler.GetHistory).Methods("GET")

	return r
}
