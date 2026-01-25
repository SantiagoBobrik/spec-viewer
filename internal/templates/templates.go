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

	layout := "templates/layouts/base.html"
	baseTmpl, err := template.ParseFS(web.Files, layout)
	if err != nil {
		log.Fatalf("Error parsing base layout: %v", err)
	}

	// Parse components (partials) so they are available to all pages
	components, err := fs.Glob(web.Files, "templates/components/*.html")
	if err != nil {
		log.Fatalf("Error globbing components: %v", err)
	}

	if len(components) > 0 {
		// ParseFiles adds the parsed templates to the existing template set
		baseTmpl, err = baseTmpl.ParseFS(web.Files, "templates/components/*.html")
		if err != nil {
			log.Fatalf("Error parsing components: %v", err)
		}
	}

	// Find all page templates (e.g., web/templates/*.html)
	// Note: in embedded FS, generic glob patterns work on forward slashes
	pages, err := fs.Glob(web.Files, "templates/*.html")
	if err != nil {
		log.Fatalf("Error globbing templates: %v", err)
	}

	for _, page := range pages {
		name := path.Base(page)
		// Extract the template name without extension (e.g., "login" from "login.html")
		tmplName := strings.TrimSuffix(name, path.Ext(name))

		// Clone the base template which now includes the components
		ts, err := baseTmpl.Clone()
		if err != nil {
			log.Fatalf("Error cloning base template for %s: %v", name, err)
		}

		// Parse the individual page template into the cloned set.
		// template.Must ensures that we panic on startup if templates are invalid.
		// Ts.ParseFS returns (*Template, error), which fits Must.
		cache[tmplName] = template.Must(ts.ParseFS(web.Files, page))
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
	err = ts.ExecuteTemplate(w, "base.html", pageData)
	if err != nil {
		log.Printf("Error executing template %s: %v", page, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
