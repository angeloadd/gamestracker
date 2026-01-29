package server

import (
	"context"
	"github.com/angeloadd/gamestracker/internal/config"
	"github.com/angeloadd/gamestracker/internal/render"
	"io/fs"
	"log/slog"
	"net/http"
)

func NewServer(
	ctx context.Context,
	cfg config.Config,
	log *slog.Logger,
	publicFS fs.FS,
	renderer *render.View,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		ctx,
		mux,
		cfg,
		log,
		publicFS,
		renderer,
	)

	var handler http.Handler = mux
	handler = logMiddleware(log)(handler)

	return handler
}
