package handler

import (
	"encoding/json"
	"net/http"
	e "notifier/internal/entity"
	s "notifier/internal/service"

	"github.com/gorilla/mux"
)

type NotifierHandler struct {
	svc *s.NotifierService
}

func NewNotifierHandler(svc *s.NotifierService) *NotifierHandler {
	return &NotifierHandler{
		svc: svc,
	}
}

func (h *NotifierHandler) CreateNotify(w http.ResponseWriter, r *http.Request) {
	var notify e.Notification
	if err := json.NewDecoder(r.Body).Decode(&notify); err != nil {
		http.Error(w, "notify cannot create", http.StatusBadRequest)
	}
	h.svc.NewNotification(notify)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Notify succesfully create ID": notify.ID,
	})
}

func (h *NotifierHandler) GetNotifyStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	status, err := h.svc.GetStatus(id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		id: status,
	})
}

func (h *NotifierHandler) DeleteNotify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.svc.DeleteNotify(id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		id: "deleted",
	})
}
