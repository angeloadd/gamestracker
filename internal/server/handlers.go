package server

import (
	"github.com/angeloadd/gamestracker/internal/render"
	"log/slog"
	"net/http"
)

func handleHome(renderer *render.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		td := &render.TemplateData{
			StringMap: map[string]string{
				"title": "Home",
				"name":  "Angelo",
			},
		}

		err := renderer.Template(w, r, "home.page.gohtml", td)
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
