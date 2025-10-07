package router

import (
	h "notifier/internal/handler"

	"github.com/gorilla/mux"
)

// – POST /notify — создание уведомлений с датой и временем отправки;
// – GET /notify/{id} — получение статуса уведомления;
// – DELETE /notify/{id} — отмена запланированного уведомления.

func NewRouter(handler h.NotifierHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/notify", handler.CreateNotify).Methods("POST")
	r.HandleFunc("/notify/{id}", handler.GetNotifyStatus).Methods("GET")
	r.HandleFunc("/notify/{id}", handler.DeleteNotify).Methods("DELETE")

	return r
}
