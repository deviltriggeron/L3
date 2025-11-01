package router

import (
	"github.com/gorilla/mux"

	"sales-tracker/internal/handler"
	"sales-tracker/internal/middleware"
)

/*
– CRUD-операции (например, финансовые транзакции или продажи);
– POST /items;
– GET /items;
– PUT /items/{id};
– DELETE /items/{id}
*/

func NewRouter(h *handler.TrackerHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.MiddlewareLogging)

	r.HandleFunc("/", h.Index)
	r.HandleFunc("/items", h.InsertItem).Methods("POST")
	r.HandleFunc("/items", h.GetAll).Methods("GET")
	r.HandleFunc("/items/{id}", h.GetItem).Methods("GET")
	r.HandleFunc("/items/{id}", h.UpdateItem).Methods("PUT")
	r.HandleFunc("/items/{id}", h.DeleteItem).Methods("DELETE")

	r.HandleFunc("/analytics", h.GetAnalytics).Methods("GET")

	return r
}
