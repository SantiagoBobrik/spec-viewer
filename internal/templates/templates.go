package templates

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/SantiagoBobrik/spec-viewer/internal/spec"
	"github.com/SantiagoBobrik/spec-viewer/web"
)

// cache holds the compiled templates for each page.
var cache = make(map[string]*template.Template)
var specFolder string

// Init parses all templates and sets the spec folder.
func Init(folder string) {
	specFolder = folder
	cache = make(map[string]*template.Template)

	baseTmpl, err := template.ParseFS(web.Files, "templates/layouts/base.html")
	if err != nil {
		log.Fatalf("Error parsing base layout: %v", err)
	}

	baseTmpl, err = baseTmpl.ParseFS(web.Files, "templates/components/*.html")
	if err != nil {
		log.Fatalf("Error parsing components: %v", err)
	}

	pages, err := fs.Glob(web.Files, "templates/*.html")
	if err != nil {
		log.Fatalf("Error globbing templates: %v", err)
	}
	if len(pages) == 0 {
		log.Fatalf("No page templates found under templates/*.html")
	}

	for _, page := range pages {
		tmplName := strings.TrimSuffix(path.Base(page), path.Ext(page))

		ts, err := baseTmpl.Clone()
		if err != nil {
			log.Fatalf("Error cloning base template for %s: %v", page, err)
		}

		ts, err = ts.ParseFS(web.Files, page)
		if err != nil {
			log.Fatalf("Error parsing page template %s: %v", page, err)
		}

		cache[tmplName] = ts
	}
}

// PageData wraps the content data with global layout data like Specs.
type PageData struct {
	Data  any
	Specs []spec.Spec
}

// Render executes the cached template.
func Render(w http.ResponseWriter, page string, data any, activePath ...string) {
	ts, ok := cache[page]
	if !ok {
		log.Printf("Template %s not found in cache", page)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	specs, err := spec.GetAll(specFolder)
	if err != nil {
		log.Printf("Error fetching specs: %v", err)
	}

	if len(activePath) > 0 {
		spec.MarkActive(specs, activePath[0])
	}

	pageData := PageData{
		Data:  data,
		Specs: specs,
	}

	// Execute the "base.html" template.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = ts.ExecuteTemplate(w, "base.html", pageData)
	if err != nil {
		log.Printf("Error executing template %s: %v", page, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
