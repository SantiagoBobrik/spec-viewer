package handlers

import (
	"net/http"
	"spec-viewer/internal/templates"
)

type HomeData struct {
	Title string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "home", HomeData{
		Title: "Home",
	})
}
