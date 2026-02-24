package main

import (
	"fmt"
	"net/http"
	"os"

	"blogging-platform-api/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/godotenv/godotenv"
)

func main() {
	// creating NewRouter

	godotenv.Load()

	conf, err := config.Load()
	if err != nil {
		os.Exit(1)
	}
	_ = conf

	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	if err := http.ListenAndServe(":2000", r); err != nil {
		fmt.Println("error ListenAndServe")
	}

	fmt.Println("server is runnign on port 8000")
}
