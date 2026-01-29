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

		w.WriteHeader(200)
	}
}

func handleHealthz(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("healthz endpoint served")
		w.WriteHeader(http.StatusNoContent)
	}
}
