package server

import (
	"context"
	"github.com/angeloadd/gamestracker/internal/config"
	"github.com/angeloadd/gamestracker/internal/render"
	"log/slog"
	"net/http"
)

func NewServer(
	ctx context.Context,
	cfg config.Config,
	log *slog.Logger,
	renderer *render.Renderer,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		ctx,
		mux,
		cfg,
		log,
		renderer,
	)

	var handler http.Handler = mux
	handler = logMiddleware(log)(handler)

	return handler
}
