package handler

import (
	"encoding/json"
	"net/http"

	"sales-tracker/internal/domain"
	"sales-tracker/internal/usecase"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TrackerHandler struct {
	svc usecase.TrackerServiceExtended
}

func NewTrackHandler(svc usecase.TrackerServiceExtended) *TrackerHandler {
	return &TrackerHandler{
		svc: svc,
	}
}

func (h *TrackerHandler) InsertItem(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.svc.Insert(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"item ID": id,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TrackerHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringID := vars["id"]

	item, err := h.svc.Get(stringID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"Item ID":           item.ID,
		"Item type":         item.Type,
		"Item amount":       item.Amount,
		"Item category":     item.Category,
		"Item description":  item.Description,
		"Item created date": item.Date,
	}

	data, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TrackerHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	vars := mux.Vars(r)
	itemID := vars["id"]

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.svc.Update(itemID, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("item succesfully updated"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TrackerHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID := vars["id"]

	err := h.svc.Delete(itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("item succesfully deleted"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TrackerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	category := q.Get("category")
	typ := q.Get("type")
	limit := q.Get("limit")
	offset := q.Get("offset")

	items, err := h.svc.GetAll(from, to, category, typ, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TrackerHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	category := q.Get("category")
	typ := q.Get("type")

	analytics, err := h.svc.GetAnalytics(from, to, category, typ)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"Sum":        analytics.Sum,
		"Average":    analytics.Avg,
		"Count":      analytics.Count,
		"Median":     analytics.Median,
		"Percentile": analytics.Percentile90,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TrackerHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
