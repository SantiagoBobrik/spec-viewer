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
	Specs []Spec
}

type Children struct {
	Name string
	Path string
}

type Spec struct {
	Name     string
	Path     string
	Children []Children
}

func ListSpecsHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all specs names
		specs, err := extractSpecs(folder)
		if err != nil {
			handleError(w, err)
			return
		}
		for i := range specs {
			err = filepath.Walk(specs[i].Path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
					relPath, err := filepath.Rel(folder, path)
					if err != nil {
						return err
					}
					specs[i].Children = append(specs[i].Children, Children{
						Name: info.Name(),
						Path: relPath,
					})
				}
				return nil
			})
			if err != nil {
				handleError(w, err)
				return
			}
		}

		templates.Render(w, "list", ListSpecsData{
			Specs: specs,
		})
	}
}

func handleError(w http.ResponseWriter, err error) {
	logger.Error("Failed to list files", "error", err)
	http.Error(w, "Failed to list files", http.StatusInternalServerError)
}

func extractSpecs(folder string) ([]Spec, error) {
	var specs []Spec
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			specs = append(specs, Spec{
				Name: info.Name(),
				Path: path,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// remove first element - root folder
	return specs[1:], nil
}
