package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SantiagoBobrik/spec-viewer/internal/templates"
	"github.com/SantiagoBobrik/spec-viewer/pkg/logger"
)

type ListSpecsData struct {
	Files []string
}

func ListSpecsHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var files []string
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				relPath, err := filepath.Rel(folder, path)
				if err != nil {
					return err
				}
				files = append(files, relPath)
			}
			return nil
		})

		if err != nil {
			logger.Error("Failed to list files", "error", err)
			http.Error(w, "Failed to list files", http.StatusInternalServerError)
			return
		}

		templates.Render(w, "list", ListSpecsData{
			Files: files,
		})
	}
}
