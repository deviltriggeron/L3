package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"imageprocessor/internal/domain"
	"imageprocessor/internal/usecase"
)

type ImageProcHandler struct {
	svc usecase.ImageProcService
}

func NewImageProcHandler(svc usecase.ImageProcService) *ImageProcHandler {
	return &ImageProcHandler{
		svc: svc,
	}
}

func (h *ImageProcHandler) Upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file not provided", http.StatusBadRequest)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			http.Error(w, fmt.Sprintf("error file close: %v", err), http.StatusInternalServerError)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "cannot read file", http.StatusInternalServerError)
		return
	}

	var options []domain.ImageProcessOption
	if r.FormValue("resize") == "true" {
		options = append(options, domain.ProcessResize)
	}
	if r.FormValue("thumbnail") == "true" {
		options = append(options, domain.ProcessThumbnail)
	}
	if r.FormValue("watermark") == "true" {
		options = append(options, domain.ProcessWatermark)
	}

	fileData := domain.FileDataType{
		FileName:    header.Filename,
		Data:        data,
		ContentType: header.Header.Get("Content-Type"),
	}
	id, err := h.svc.Upload(fileData, options...)
	if err != nil {
		http.Error(w, fmt.Sprintf("upload failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(
		map[string]string{
			"File upload": id,
		}); err != nil {
		http.Error(w, fmt.Sprintf("upload failed: %v", err), http.StatusInternalServerError)
	}
}

func (h *ImageProcHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	option := r.URL.Query().Get("variant")
	if option == "" {
		option = "resized"
	}

	data, contentType, err := h.svc.Get(id, option)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get image: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		http.Error(w, fmt.Sprintf("error write data: %v", err), http.StatusBadRequest)
	}
}

func (h *ImageProcHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.svc.DeleteImage(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error delete image: %v", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(
		map[string]string{
			"File delete": id,
		}); err != nil {
		http.Error(w, fmt.Sprintf("delete failed: %v", err), http.StatusInternalServerError)
	}
}

func (h *ImageProcHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
