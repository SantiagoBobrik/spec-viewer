package templates

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// cache holds the compiled templates for each page.
var cache = make(map[string]*template.Template)

func init() {
	layout := filepath.Join("web", "templates", "layouts", "base.html")
	baseTmpl, err := template.ParseFiles(layout)
	if err != nil {
		log.Fatalf("Error parsing base layout: %v", err)
	}

	// Parse components (partials) so they are available to all pages
	components, err := filepath.Glob(filepath.Join("web", "templates", "components", "*.html"))
	if err != nil {
		log.Fatalf("Error globbing components: %v", err)
	}

	if len(components) > 0 {
		// ParseFiles adds the parsed templates to the existing template set
		baseTmpl, err = baseTmpl.ParseFiles(components...)
		if err != nil {
			log.Fatalf("Error parsing components: %v", err)
		}
	}

	// Find all page templates (e.g., web/templates/*.html)
	pages, err := filepath.Glob(filepath.Join("web", "templates", "*.html"))
	if err != nil {
		log.Fatalf("Error globbing templates: %v", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// Extract the template name without extension (e.g., "login" from "login.html")
		tmplName := strings.TrimSuffix(name, filepath.Ext(name))

		// Clone the base template which now includes the components
		ts, err := baseTmpl.Clone()
		if err != nil {
			log.Fatalf("Error cloning base template for %s: %v", name, err)
		}

		// Parse the individual page template into the cloned set.
		// template.Must ensures that we panic on startup if templates are invalid.
		cache[tmplName] = template.Must(ts.ParseFiles(page))
	}
}

// Render executes the cached template.
func Render(w http.ResponseWriter, page string, data any) {
	ts, ok := cache[page]
	if !ok {
		log.Printf("Template %s not found in cache", page)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the "base.html" template.
	err := ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Error executing template %s: %v", page, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
