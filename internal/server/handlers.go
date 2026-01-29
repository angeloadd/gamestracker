package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

func handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, World!")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleHealth(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.InfoContext(r.Context(), "healthz endpoint served")
		w.WriteHeader(http.StatusNoContent)
	}
}
