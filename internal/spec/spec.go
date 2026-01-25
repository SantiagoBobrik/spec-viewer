package spec

import (
	"os"
	"path/filepath"
	"strings"
)

type Children struct {
	Name string
	Path string
}

type Spec struct {
	Name     string
	Path     string
	Children []Children
}

func GetAll(folder string) ([]Spec, error) {
	var specs []Spec
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip root folder itself if it matches 'folder' exactly, or handle as needed.
			if path == folder {
				return nil
			}

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

	// Now populate children for each spec
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
			return nil, err
		}
	}

	return specs, nil
}
