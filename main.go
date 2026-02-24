package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	// creating NewRouter
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	if err := http.ListenAndServe(":8000", r); err != nil {
		fmt.Println("error ListenAndServe")
	}
}
