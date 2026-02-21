package spec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAll_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	specs, err := GetAll(dir)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("expected 0 specs, got %d", len(specs))
	}
}

func TestGetAll_OnlyMarkdownFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "spec1.md", "# Spec 1")
	writeFile(t, dir, "spec2.md", "# Spec 2")
	writeFile(t, dir, "readme.txt", "not markdown")
	writeFile(t, dir, "image.png", "binary data")

	specs, err := GetAll(dir)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(specs) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(specs))
	}

	names := map[string]bool{}
	for _, s := range specs {
		names[s.Name] = true
		if s.IsDir {
			t.Errorf("spec %q should not be a directory", s.Name)
		}
	}
	if !names["spec1.md"] {
		t.Error("expected spec1.md in results")
	}
	if !names["spec2.md"] {
		t.Error("expected spec2.md in results")
	}
}

func TestGetAll_HiddenFilesExcluded(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, ".hidden.md", "# Hidden")
	writeFile(t, dir, "visible.md", "# Visible")
	mkdir(t, dir, ".hidden-dir")
	writeFile(t, filepath.Join(dir, ".hidden-dir"), "secret.md", "# Secret")

	specs, err := GetAll(dir)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}
	if specs[0].Name != "visible.md" {
		t.Errorf("expected visible.md, got %s", specs[0].Name)
	}
}

func TestGetAll_RecursiveDirectoryScanning(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "root.md", "# Root")
	mkdir(t, dir, "subdir")
	writeFile(t, filepath.Join(dir, "subdir"), "nested.md", "# Nested")
	mkdir(t, dir, "subdir", "deep")
	writeFile(t, filepath.Join(dir, "subdir", "deep"), "deep.md", "# Deep")

	specs, err := GetAll(dir)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(specs) != 2 {
		t.Fatalf("expected 2 top-level specs (root.md and subdir), got %d", len(specs))
	}

	// Find the file and the directory
	var fileSpec, dirSpec *Spec
	for i := range specs {
		if specs[i].Name == "root.md" {
			fileSpec = &specs[i]
		}
		if specs[i].Name == "subdir" {
			dirSpec = &specs[i]
		}
	}

	if fileSpec == nil {
		t.Fatal("expected root.md in results")
	}
	if fileSpec.Path != "root.md" {
		t.Errorf("expected path 'root.md', got %q", fileSpec.Path)
	}
	if fileSpec.IsDir {
		t.Error("root.md should not be a directory")
	}

	if dirSpec == nil {
		t.Fatal("expected subdir in results")
	}
	if !dirSpec.IsDir {
		t.Error("subdir should be a directory")
	}
	if dirSpec.Path != "subdir" {
		t.Errorf("expected path 'subdir', got %q", dirSpec.Path)
	}

	// Check children
	if len(dirSpec.Children) != 2 {
		t.Fatalf("expected 2 children in subdir (deep and nested.md), got %d", len(dirSpec.Children))
	}

	// Find the deep directory within children
	var deepDir *Spec
	for i := range dirSpec.Children {
		if dirSpec.Children[i].Name == "deep" {
			deepDir = &dirSpec.Children[i]
		}
	}
	if deepDir == nil {
		t.Fatal("expected deep directory in subdir children")
	}
	if len(deepDir.Children) != 1 {
		t.Fatalf("expected 1 child in deep directory, got %d", len(deepDir.Children))
	}
	if deepDir.Children[0].Name != "deep.md" {
		t.Errorf("expected deep.md in deep dir, got %s", deepDir.Children[0].Name)
	}
	if deepDir.Children[0].Path != filepath.Join("subdir", "deep", "deep.md") {
		t.Errorf("expected nested path, got %q", deepDir.Children[0].Path)
	}
}

func TestGetAll_EmptySubdirectory(t *testing.T) {
	dir := t.TempDir()
	mkdir(t, dir, "empty")

	specs, err := GetAll(dir)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec (empty dir), got %d", len(specs))
	}
	if !specs[0].IsDir {
		t.Error("expected empty to be a directory")
	}
	if len(specs[0].Children) != 0 {
		t.Errorf("expected 0 children in empty dir, got %d", len(specs[0].Children))
	}
}

func TestGetAll_NonExistentDirectory(t *testing.T) {
	_, err := GetAll("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
}

func TestGetAll_DirectoryWithOnlyNonMarkdownFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "readme.txt", "text")
	writeFile(t, dir, "data.json", "{}")
	writeFile(t, dir, "script.sh", "#!/bin/bash")

	specs, err := GetAll(dir)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("expected 0 specs (no .md files), got %d", len(specs))
	}
}

func TestGetAll_RelativePaths(t *testing.T) {
	dir := t.TempDir()
	mkdir(t, dir, "docs")
	writeFile(t, filepath.Join(dir, "docs"), "api.md", "# API")

	specs, err := GetAll(dir)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}

	docsDir := specs[0]
	if docsDir.Path != "docs" {
		t.Errorf("expected path 'docs', got %q", docsDir.Path)
	}
	if len(docsDir.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(docsDir.Children))
	}
	if docsDir.Children[0].Path != filepath.Join("docs", "api.md") {
		t.Errorf("expected path 'docs/api.md', got %q", docsDir.Children[0].Path)
	}
}

func TestMarkActive_TopLevelFile(t *testing.T) {
	specs := []Spec{
		{Name: "a.md", Path: "a.md"},
		{Name: "b.md", Path: "b.md"},
		{Name: "c.md", Path: "c.md"},
	}

	MarkActive(specs, "b.md")

	if specs[0].Active {
		t.Error("a.md should not be active")
	}
	if !specs[1].Active {
		t.Error("b.md should be active")
	}
	if specs[2].Active {
		t.Error("c.md should not be active")
	}
}

func TestMarkActive_NestedFile(t *testing.T) {
	specs := []Spec{
		{Name: "root.md", Path: "root.md"},
		{
			Name:  "subdir",
			Path:  "subdir",
			IsDir: true,
			Children: []Spec{
				{Name: "nested.md", Path: "subdir/nested.md"},
				{Name: "other.md", Path: "subdir/other.md"},
			},
		},
	}

	MarkActive(specs, "subdir/nested.md")

	if specs[0].Active {
		t.Error("root.md should not be active")
	}
	if specs[1].Active {
		t.Error("subdir itself should not be active")
	}
	if !specs[1].Children[0].Active {
		t.Error("subdir/nested.md should be active")
	}
	if specs[1].Children[1].Active {
		t.Error("subdir/other.md should not be active")
	}
}

func TestMarkActive_NonExistentPath(t *testing.T) {
	specs := []Spec{
		{Name: "a.md", Path: "a.md"},
		{Name: "b.md", Path: "b.md"},
	}

	MarkActive(specs, "nonexistent.md")

	for _, s := range specs {
		if s.Active {
			t.Errorf("%s should not be active", s.Name)
		}
	}
}

func TestMarkActive_EmptySpecs(t *testing.T) {
	var specs []Spec
	// Should not panic
	MarkActive(specs, "anything.md")
}

func TestMarkActive_EmptyActivePath(t *testing.T) {
	specs := []Spec{
		{Name: "a.md", Path: "a.md"},
	}

	MarkActive(specs, "")

	if specs[0].Active {
		t.Error("a.md should not be active when activePath is empty")
	}
}

func TestMarkActive_DeeplyNested(t *testing.T) {
	specs := []Spec{
		{
			Name:  "level1",
			Path:  "level1",
			IsDir: true,
			Children: []Spec{
				{
					Name:  "level2",
					Path:  "level1/level2",
					IsDir: true,
					Children: []Spec{
						{Name: "deep.md", Path: "level1/level2/deep.md"},
					},
				},
			},
		},
	}

	MarkActive(specs, "level1/level2/deep.md")

	if specs[0].Active {
		t.Error("level1 should not be active")
	}
	if specs[0].Children[0].Active {
		t.Error("level1/level2 should not be active")
	}
	if !specs[0].Children[0].Children[0].Active {
		t.Error("level1/level2/deep.md should be active")
	}
}

// Helper functions

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write file %s: %v", name, err)
	}
}

func mkdir(t *testing.T, parts ...string) {
	t.Helper()
	path := filepath.Join(parts...)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatalf("failed to create directory %s: %v", path, err)
	}
}
