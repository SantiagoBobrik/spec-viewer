package spec

import (
	"os"
	"path/filepath"
	"strings"
)

type Spec struct {
	Name     string
	Path     string
	IsDir    bool
	Active   bool
	Children []Spec
}

func GetAll(root string) ([]Spec, error) {
	return scanDir(root, "")
}

func scanDir(root string, relBase string) ([]Spec, error) {
	fullPath := filepath.Join(root, relBase)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var specs []Spec
	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files/dirs
		if strings.HasPrefix(name, ".") {
			continue
		}

		// Skip non-markdown files if it's a file
		if !entry.IsDir() && !strings.HasSuffix(name, ".md") {
			continue
		}

		relPath := filepath.Join(relBase, name)
		item := Spec{
			Name:  name,
			Path:  relPath,
			IsDir: entry.IsDir(),
		}

		if entry.IsDir() {
			children, err := scanDir(root, relPath)
			if err != nil {
				return nil, err
			}
			item.Children = children
		}

		specs = append(specs, item)
	}

	return specs, nil
}

func MarkActive(specs []Spec, activePath string) {
	for i := range specs {
		if specs[i].Path == activePath {
			specs[i].Active = true
		}
		if specs[i].IsDir {
			MarkActive(specs[i].Children, activePath)
		}
	}
}
