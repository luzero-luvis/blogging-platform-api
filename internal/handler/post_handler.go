package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"blogging-platform-api/internal/model"
	"blogging-platform-api/internal/service"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("invalide json", "error", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	post, err := h.service.Create(&req)
	if err != nil {
		slog.Error("failed to create post", "error", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	slog.Info("post created for id", "id", post.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(post)
}
