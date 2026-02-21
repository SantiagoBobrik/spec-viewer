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
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// TOCEntry represents a single heading in the table of contents.
type TOCEntry struct {
	Level int
	Text  string
	ID    string
}

type ViewerData struct {
	Title   string
	Content template.HTML
	TOC     []TOCEntry
}

// md is the shared Goldmark instance configured with auto heading IDs for TOC generation.
var md = goldmark.New(
	goldmark.WithExtensions(
		extension.Table,
		extension.Strikethrough,
		extension.Linkify,
		extension.TaskList,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

// renderMarkdown validates the file parameter, reads the markdown file, and
// converts it to HTML. It returns the cleaned path, the rendered HTML bytes,
// and the TOC entries. If an error occurs, it writes an appropriate HTTP
// response and returns false.
func renderMarkdown(folder string, w http.ResponseWriter, r *http.Request) (string, []byte, []TOCEntry, bool) {
	fileParam := r.URL.Query().Get("file")
	if fileParam == "" {
		logger.Info("File not specified - redirecting to home")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return "", nil, nil, false
	}

	// Security check: prevent directory traversal
	cleanPath := filepath.Clean(fileParam)
	if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, "/") {
		logger.Info("Invalid file path - redirecting to home")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return "", nil, nil, false
	}

	fullPath := filepath.Join(folder, cleanPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		logger.Error("Failed to read file", "file", fullPath, "error", err)
		if os.IsNotExist(err) {
			logger.Info("File not found - redirecting to home")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return "", nil, nil, false
		}
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return "", nil, nil, false
	}

	// Parse markdown into AST and extract TOC entries.
	reader := text.NewReader(content)
	doc := md.Parser().Parse(reader)
	toc := extractTOC(doc, content)

	// Render markdown to HTML.
	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, content, doc); err != nil {
		logger.Error("Failed to render markdown", "error", err)
		http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
		return "", nil, nil, false
	}

	return cleanPath, buf.Bytes(), toc, true
}

func ViewSpecHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cleanPath, html, toc, ok := renderMarkdown(folder, w, r)
		if !ok {
			return
		}

		templates.Render(w, "viewer", ViewerData{
			Title:   cleanPath,
			Content: template.HTML(html),
			TOC:     toc,
		}, cleanPath)
	}
}

// ViewContentHandler returns only the rendered markdown HTML fragment,
// without the full page template wrapper. This is used by the WebSocket
// client to update content in-place without a full page reload.
func ViewContentHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, html, _, ok := renderMarkdown(folder, w, r)
		if !ok {
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(html)
	}
}

// extractTOC walks the Goldmark AST and collects heading entries for the table of contents.
func extractTOC(doc ast.Node, source []byte) []TOCEntry {
	var entries []TOCEntry

	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		heading, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		// Extract the heading text from child nodes.
		var textBuf bytes.Buffer
		for child := heading.FirstChild(); child != nil; child = child.NextSibling() {
			if t, ok := child.(*ast.Text); ok {
				textBuf.Write(t.Segment.Value(source))
			} else {
				// For non-text children (e.g., code spans, emphasis), collect their text content.
				_ = ast.Walk(child, func(cn ast.Node, entering bool) (ast.WalkStatus, error) {
					if entering {
						if ct, ok := cn.(*ast.Text); ok {
							textBuf.Write(ct.Segment.Value(source))
						}
					}
					return ast.WalkContinue, nil
				})
			}
		}

		// Get the auto-generated heading ID.
		id, found := heading.AttributeString("id")
		if !found {
			return ast.WalkContinue, nil
		}

		idStr := ""
		switch v := id.(type) {
		case []byte:
			idStr = string(v)
		case string:
			idStr = v
		}

		entries = append(entries, TOCEntry{
			Level: heading.Level,
			Text:  textBuf.String(),
			ID:    idStr,
		})

		return ast.WalkContinue, nil
	})

	return entries
}
