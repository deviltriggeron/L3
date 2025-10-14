package handler

import (
	"encoding/json"
	"net/http"
	svc "shortener/internal/service"

	"github.com/gorilla/mux"
)

type ShortenerHandler struct {
	svc svc.ShortenerService
}

func NewShortenerHandler(svc svc.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{
		svc: svc,
	}
}

func (h *ShortenerHandler) NewShorten(w http.ResponseWriter, r *http.Request) {
	var url string
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	shortUrl, err := h.svc.NewShorten(r.Context(), url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Shortest url succesfully create": shortUrl,
	})
}

func (h *ShortenerHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["short_url"]

	longURL, err := h.svc.Redirect(r.Context(), shortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.svc.LogTransition(shortUrl, r.UserAgent(), r.RemoteAddr)
	http.Redirect(w, r, longURL, http.StatusFound)
}

func (h *ShortenerHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["short_url"]

	transitions, err := h.svc.GetAnalyticsData(r.Context(), shortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	resp := map[string]interface{}{
		"analytics": transitions,
		"count":     len(transitions),
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
