package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"warehouse-control/internal/domain"
	"warehouse-control/internal/middleware"
	"warehouse-control/internal/usecase"
)

type ControllerHandler struct {
	svc usecase.ControllerService
}

func NewHandler(svc usecase.ControllerService) ControllerHandler {
	return ControllerHandler{
		svc: svc,
	}
}

func (h ControllerHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func (h ControllerHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r)
	if user.Role != domain.RoleAdmin {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var item domain.Item
	var err error

	if err = json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.svc.AddItem(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h ControllerHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	item, err := h.svc.GetItem(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"Item ID":          item.ID,
		"Item name":        item.Product,
		"Item price":       item.Price,
		"Item description": item.Description,
		"Item count":       item.Count,
		"Item create date": item.CreateDate,
	}

	data, err := json.MarshalIndent(resp, " ", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h ControllerHandler) GetAllItem(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.GetAllItem()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

func (h ControllerHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r)
	if user.Role != domain.RoleAdmin && user.Role != domain.RoleManager {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	var item domain.Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.svc.UpdateItem(idStr, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("Item succesfully update"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h ControllerHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r)
	if user.Role != domain.RoleAdmin {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	err := h.svc.DeleteItem(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("Item succesfully deleted"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
