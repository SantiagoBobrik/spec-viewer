package handlers

import (
	"net/http"

	"github.com/SantiagoBobrik/spec-viewer/internal/templates"
)

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		templates.Render(w, "404", nil)
	}
}
