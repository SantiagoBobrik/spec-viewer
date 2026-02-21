package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/SantiagoBobrik/spec-viewer/internal/templates"
	"github.com/SantiagoBobrik/spec-viewer/pkg/logger"
)

// testSpecDir is a temporary directory used as the spec folder for templates
// and the ViewSpecHandler during tests.
var testSpecDir string

func TestMain(m *testing.M) {
	// Silence logger output.
	logger.SetOutput(io.Discard)

	// Create a temp directory to act as the spec folder.
	dir, err := os.MkdirTemp("", "handlers-test-specs-*")
	if err != nil {
		panic(err)
	}
	testSpecDir = dir

	// Write a sample markdown file to the spec folder.
	err = os.WriteFile(filepath.Join(dir, "sample.md"), []byte("# Sample\n\nHello world"), 0644)
	if err != nil {
		panic(err)
	}

	// Initialize the template cache with the spec folder.
	templates.Init(dir)

	code := m.Run()

	_ = os.RemoveAll(dir)
	os.Exit(code)
}

// --- HomeHandler tests ---

func TestHomeHandler_ReturnsOK(t *testing.T) {
	handler := HomeHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type 'text/html; charset=utf-8', got %q", contentType)
	}

	body := rr.Body.String()
	if len(body) == 0 {
		t.Error("expected non-empty response body")
	}
}

// --- NotFoundHandler tests ---

func TestNotFoundHandler_Returns404(t *testing.T) {
	handler := NotFoundHandler()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNotFoundHandler_ReturnsHTML(t *testing.T) {
	handler := NotFoundHandler()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type 'text/html; charset=utf-8', got %q", contentType)
	}
}

// --- ViewSpecHandler tests ---

func TestViewSpecHandler_NoFileParam_Redirects(t *testing.T) {
	handler := ViewSpecHandler(testSpecDir)
	req := httptest.NewRequest(http.MethodGet, "/view", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}
	location := rr.Header().Get("Location")
	if location != "/" {
		t.Errorf("expected redirect to '/', got %q", location)
	}
}

func TestViewSpecHandler_EmptyFileParam_Redirects(t *testing.T) {
	handler := ViewSpecHandler(testSpecDir)
	req := httptest.NewRequest(http.MethodGet, "/view?file=", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}
}

func TestViewSpecHandler_DirectoryTraversal_Redirects(t *testing.T) {
	tests := []struct {
		name      string
		fileParam string
	}{
		{"double dot", "../../etc/passwd"},
		{"encoded traversal", "../secret.md"},
		{"absolute path", "/etc/passwd"},
		{"absolute path with md", "/tmp/malicious.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ViewSpecHandler(testSpecDir)
			req := httptest.NewRequest(http.MethodGet, "/view?file="+tt.fileParam, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusSeeOther {
				t.Errorf("expected redirect (303) for file=%q, got %d", tt.fileParam, rr.Code)
			}
			location := rr.Header().Get("Location")
			if location != "/" {
				t.Errorf("expected redirect to '/', got %q", location)
			}
		})
	}
}

func TestViewSpecHandler_MissingFile_Redirects(t *testing.T) {
	handler := ViewSpecHandler(testSpecDir)
	req := httptest.NewRequest(http.MethodGet, "/view?file=nonexistent.md", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d for missing file, got %d", http.StatusSeeOther, rr.Code)
	}
}

func TestViewSpecHandler_ValidFile_ReturnsOK(t *testing.T) {
	handler := ViewSpecHandler(testSpecDir)
	req := httptest.NewRequest(http.MethodGet, "/view?file=sample.md", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type 'text/html; charset=utf-8', got %q", contentType)
	}

	body := rr.Body.String()
	if len(body) == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestViewSpecHandler_ValidFile_RendersMarkdown(t *testing.T) {
	handler := ViewSpecHandler(testSpecDir)
	req := httptest.NewRequest(http.MethodGet, "/view?file=sample.md", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	// The markdown "# Sample" should be rendered as an <h1> tag with an auto-generated id.
	if !containsSubstring(body, "<h1 id=\"sample\">Sample</h1>") {
		t.Error("expected rendered HTML to contain '<h1 id=\"sample\">Sample</h1>'")
	}
	// The paragraph text should appear.
	if !containsSubstring(body, "Hello world") {
		t.Error("expected rendered HTML to contain 'Hello world'")
	}
}

func TestViewSpecHandler_NestedFile_ReturnsOK(t *testing.T) {
	// Create a nested directory structure in the spec folder.
	subdir := filepath.Join(testSpecDir, "nested")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "deep.md"), []byte("# Deep"), 0644); err != nil {
		t.Fatalf("failed to write deep.md: %v", err)
	}
	defer func() { _ = os.RemoveAll(subdir) }()

	handler := ViewSpecHandler(testSpecDir)
	req := httptest.NewRequest(http.MethodGet, "/view?file=nested/deep.md", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && contains(s, substr))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
