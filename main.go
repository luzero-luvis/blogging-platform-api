package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"blogging-platform-api/internal/config"
	"blogging-platform-api/internal/database"
	"blogging-platform-api/internal/logger"
	"blogging-platform-api/internal/middleware"

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
	_ = conf

	// connnect to database
	db, err := database.Connect(conf.DBURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
	} else {
		defer db.Close()
		slog.Info("connected to database")
	}

	r := chi.NewRouter()

	// this is middleware
	r.Use(middleware.Loggingmiddleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("postgres is unavailble"))
		}
		w.Write([]byte("ok"))
	})

	// port

	port := conf.PORT
	if port == "" {
		port = "8080"
	}

	slog.Info("server startin at port", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("server failed to start", "error", err)
	}

	fmt.Println("server is runnign on port 8000")
}
