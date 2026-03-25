package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"blogging-platform-api/internal/model"
	"blogging-platform-api/internal/service"

	"github.com/go-chi/chi/v5"
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	post, err := h.service.Create(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAll()
	if err != nil {
		slog.Error("failed to get data", "error", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get post"})
		return
	}
	if posts == nil {
		posts = []model.Post{}
	}
	slog.Info("fetched all data", "respose", "sucess")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("invalid id", "id", idStr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})

	}

	post, err := h.service.GetById(id)
	if err != nil {
		slog.Warn("id doesn't exists", "warning", "there is no id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string{"warning": "id doesn't exist"})
		return
	}

	slog.Info("data for", "id", post.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(post)
}
