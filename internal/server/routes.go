package server

import (
	"context"
	"github.com/angeloadd/gamestracker/internal/config"
	"github.com/angeloadd/gamestracker/internal/render"
	"log/slog"
	"net/http"
)

func addRoutes(
	_ context.Context,
	mux *http.ServeMux,
	_ config.Config,
	log *slog.Logger,
	renderer *render.Renderer,
) {
	mux.HandleFunc("GET /{$}", handleHome(renderer))
	mux.HandleFunc("GET /healthz", handleHealth(log))
	mux.Handle("/", http.NotFoundHandler())
}
