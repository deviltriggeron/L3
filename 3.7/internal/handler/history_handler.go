package handler

import (
	"encoding/json"
	"net/http"

	"warehouse-control/internal/usecase"
)

type HistoryHandler struct {
	svc usecase.HistoryService
}

func NewHistoryHandler(svc usecase.HistoryService) HistoryHandler {
	return HistoryHandler{
		svc: svc,
	}
}

func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	history, err := h.svc.GetHistory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data, err := json.MarshalIndent(history, " ", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
