package handlers

import (
	"net/http"

	"github.com/SantiagoBobrik/spec-viewer/internal/templates"
)

func ListSpecsHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates.Render(w, "home", nil)
	}
}
