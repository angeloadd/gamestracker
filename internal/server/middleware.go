package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func logMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.DebugContext(r.Context(), fmt.Sprintf("incoming request from %s: %s %s", r.RemoteAddr, r.Method, r.URL.Path))

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(w, r)
			log.Info("",
				"status", wrapped.Status(),
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		})
	}
}
