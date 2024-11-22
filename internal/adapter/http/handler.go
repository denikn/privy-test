package http

import (
	"beprivytest/internal/port"
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Code         int         `json:"code"`
	Result       interface{} `json:"result,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

type facesCountResponse struct {
	ProcessedPath string `json:"processedPath"`
	FaceCount     int    `json:"faceCount"`
}

type ImageHandler struct {
	service port.ImageService
}

func NewImageHandler(service port.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

func (h *ImageHandler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		h.respondJSON(w, http.StatusBadRequest, jsonResponse{Code: http.StatusBadRequest, ErrorMessage: "Failed to read file"})
		return
	}
	defer file.Close()

	filename := header.Filename
	_, err = h.service.SaveImage(file, filename)
	if err != nil {
		h.respondJSON(w, http.StatusInternalServerError, jsonResponse{Code: http.StatusInternalServerError, ErrorMessage: "Failed to process image"})
		return
	}

	h.respondJSON(w, http.StatusOK, jsonResponse{Code: http.StatusOK, Result: "Image uploaded successfully"})
}

func (h *ImageHandler) GetFacesCount(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	processedPath, faceCount, err := h.service.ProcessImage(filename)

	if err != nil {
		h.respondJSON(w, http.StatusInternalServerError, jsonResponse{Code: http.StatusInternalServerError, ErrorMessage: "Failed to get faces count"})
		return
	}

	response := facesCountResponse{
		ProcessedPath: processedPath,
		FaceCount:     faceCount,
	}

	h.respondJSON(w, http.StatusOK, jsonResponse{Code: http.StatusOK, Result: response})
}

func (h *ImageHandler) ServeProcessedImage(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	processedImagePath := h.service.GetProcessedImagePath(filename)
	http.ServeFile(w, r, processedImagePath)
}
