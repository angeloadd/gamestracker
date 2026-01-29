package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/angeloadd/gamestracker/internal/config"
	"github.com/angeloadd/gamestracker/internal/logging"
	"github.com/angeloadd/gamestracker/internal/render"
	"github.com/angeloadd/gamestracker/internal/server"
	"github.com/angeloadd/gamestracker/web"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.LookupEnv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	_ []string,
	lookenv func(string) (string, bool),
	_ io.Reader,
	_, stderr io.Writer,
) error {
	cfg := config.NewConfig(lookenv)
	log := logging.NewLogger(cfg.Logs, stderr)
	renderer := render.NewView(cfg, log)

	if !cfg.App.Debug {
		cache, err := renderer.CreateTemplateCache()
		if err != nil {
			return fmt.Errorf("error creating renderer: %v", err)
		}

		renderer.TemplateCache = cache
	}

	srv := server.NewServer(ctx, cfg, log, web.PublicFS, renderer)
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelError),
	}

	go func() {
		log.Debug(fmt.Sprintf("listening on %s\n", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(fmt.Sprintf("error listening and serving: %s\n", err))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error(fmt.Sprintf("error shutting down http server: %s\n", err))
		}
	}()
	wg.Wait()
	return nil
}
