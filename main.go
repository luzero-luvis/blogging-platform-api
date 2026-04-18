package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"blogging-platform-api/internal/config"
	"blogging-platform-api/internal/database"
	"blogging-platform-api/internal/handler"
	"blogging-platform-api/internal/logger"
	"blogging-platform-api/internal/middleware"
	"blogging-platform-api/internal/repository"
	"blogging-platform-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/godotenv/godotenv"
)

func main() {
	// creating NewRouter

	godotenv.Load()

	logger.Setup(os.Getenv("ENV"))

	conf, err := config.Load()
	if err != nil {
		os.Exit(1)
	}

	// connnect to database
	db, err := database.Connect(conf.DBURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
	} else {
		defer db.Close()
		slog.Info("connected to database")
	}

	repo := repository.NewPostRepository(db)
	svc := service.NewPostService(repo)
	h := handler.NewPostHandler(svc)

	r := chi.NewRouter()

	// this is middleware
	r.Use(middleware.Loggingmiddleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if db == nil || db.Ping() != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "service is not healthy",
				"data": map[string]bool{
					"postgress": false,
				},
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "service is healthy",
			"data": map[string]bool{
				"postgress": true,
			},
		})
	})

	r.Post("/posts", h.Create)
	r.Get("/posts", h.GetAll)
	r.Get("/posts/{id}", h.GetByID)
	r.Put("/posts/{id}", h.Put)
	r.Delete("/posts/{id}", h.Del)
	// port

	port := conf.PORT
	if port == "" {
		port = "8080"
	}

	slog.Info("server startin at port", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("server failed to start", "error", err)
	}
}
