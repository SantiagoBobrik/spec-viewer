package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SantiagoBobrik/spec-viewer/internal/templates"
	"github.com/SantiagoBobrik/spec-viewer/pkg/logger"
	"github.com/yuin/goldmark"
)

type ViewerData struct {
	Title   string
	Content template.HTML
}

func ViewSpecHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileParam := r.URL.Query().Get("file")
		if fileParam == "" {
			http.Error(w, "File not specified", http.StatusBadRequest)
			return
		}

		// Security check: prevent directory traversal
		cleanPath := filepath.Clean(fileParam)
		if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, "/") {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join(folder, cleanPath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			logger.Error("Failed to read file", "file", fullPath, "error", err)
			if os.IsNotExist(err) {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		if err := goldmark.Convert(content, &buf); err != nil {
			logger.Error("Failed to render markdown", "error", err)
			http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
			return
		}

		templates.Render(w, "viewer", ViewerData{
			Title:   cleanPath,
			Content: template.HTML(buf.String()),
		}, cleanPath)
	}
}
