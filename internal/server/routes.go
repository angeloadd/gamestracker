package server

import (
	"context"
	"github.com/angeloadd/gamestracker/internal/config"
	"log/slog"
	"net/http"
)

func addRoutes(
	_ context.Context,
	mux *http.ServeMux,
	_ config.Config,
	log *slog.Logger,
) {
	mux.Handle("/", handleHome())
	mux.HandleFunc("/healthz", handleHealthz(log))
	mux.Handle("/", http.NotFoundHandler())
}
