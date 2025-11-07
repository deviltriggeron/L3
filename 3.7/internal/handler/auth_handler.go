package handler

import (
	"encoding/json"
	"net/http"

	"warehouse-control/internal/interfaces"
	"warehouse-control/internal/usecase"
)

type AuthHandler struct {
	svc usecase.AuthService
	tp  interfaces.TokenProvide
}

func NewAuthHandler(svc usecase.AuthService, tp interfaces.TokenProvide) AuthHandler {
	return AuthHandler{
		svc: svc,
		tp:  tp,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	type request struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.svc.Login(req.User, req.Password)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
