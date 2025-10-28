package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"eventbooker/internal/domain"
	"eventbooker/internal/usecase"
)

type EventBookHandler struct {
	svc usecase.EventBookService
}

func NewHandler(svc usecase.EventBookService) *EventBookHandler {
	return &EventBookHandler{
		svc: svc,
	}
}

func (h *EventBookHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var e domain.EventBook

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := uuid.New()
	e.EventID = id
	err := h.svc.CreateEvent(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"Event ID": id,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *EventBookHandler) Book(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["id"]

	var userID string
	if err := json.NewDecoder(r.Body).Decode(&userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookingID, err := h.svc.Booking(userID, eventID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if bookingID == uuid.Nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "This is a free event, you have successfully registered",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "The payment window will close in 15 minutes",
		"booking_id": bookingID,
	})
}

func (h *EventBookHandler) Payment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	var eventID string
	if err := json.NewDecoder(r.Body).Decode(&eventID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.svc.Payment(eventID, bookID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if jsonErr := json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		}); jsonErr != nil {
			http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("successfully paid"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *EventBookHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	event, err := h.svc.GetEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"Event ID":           event.EventID,
		"Event date":         event.Date,
		"Event info":         event.EventInfo,
		"Event Organizer":    event.Organizer,
		"Event seats count":  event.SeatsCount,
		"Event free":         event.ForFree,
		"Event participants": len(event.Participants),
		"Event price":        event.Price,
		"Event create date":  event.CreateDate,
	}

	data, err := json.MarshalIndent(resp, "", "  ")
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

func (h *EventBookHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var username string
	id := uuid.New()

	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := domain.User{
		UserID: id,
		Name:   username,
	}

	err := h.svc.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"user ID": id,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *EventBookHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
