package server

import (
	"context"
	"github.com/angeloadd/gamestracker/internal/config"
	"log/slog"
	"net/http"
)

func NewServer(
	ctx context.Context,
	cfg config.Config,
	log *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		ctx,
		mux,
		cfg,
		log,
	)

	var handler http.Handler = mux
	handler = logMiddleware(log)(handler)

	return handler
}
