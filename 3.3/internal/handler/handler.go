package handler

import (
	e "commentTree/internal/entity"
	"commentTree/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CommetsHandler struct {
	svc *service.CommentsService
}

func NewCommentsHandler(svc *service.CommentsService) *CommetsHandler {
	return &CommetsHandler{
		svc: svc,
	}
}

func (h *CommetsHandler) NewComments(w http.ResponseWriter, r *http.Request) {
	var comments e.CommentResponse
	if err := json.NewDecoder(r.Body).Decode(&comments); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := h.svc.Comments(r.Context(), comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	resp := map[string]interface{}{
		"Comment posted": c.Comment,
		"Comment ID":     c.CommentID,
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *CommetsHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("parent")

	comments, err := h.svc.GetComments(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	data, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *CommetsHandler) GetAllParentComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.svc.GetAllParentComments(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	data, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *CommetsHandler) DeleteComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.svc.DeleteComments(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Comment deleted": id,
	})
}

func (h *CommetsHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
