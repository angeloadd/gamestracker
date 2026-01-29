package server

import (
	"context"
	"github.com/angeloadd/gamestracker/internal/config"
	"github.com/angeloadd/gamestracker/internal/render"
	"io/fs"
	"log/slog"
	"net/http"
)

func addRoutes(
	_ context.Context,
	mux *http.ServeMux,
	_ config.Config,
	log *slog.Logger,
	publicFS fs.FS,
	renderer *render.View,
) {
	mux.HandleFunc("GET /{$}", handleHome(renderer))
	mux.HandleFunc("GET /healthz", handleHealth(log))

	subFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Error("Failed to load publicFS", err)
	}
	mux.Handle("GET /s/", http.StripPrefix("/s/", http.FileServerFS(subFS)))
	mux.Handle("/", http.NotFoundHandler())
}
