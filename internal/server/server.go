package server

import (
	"net/http"
	"strings"

	"github.com/SantiagoBobrik/spec-viewer/internal/handlers"
	"github.com/SantiagoBobrik/spec-viewer/internal/socket"

	"github.com/gorilla/mux"
)

type Config struct {
	Port   string
	Folder string
}

// noDirectoryListing intercepta peticiones a directorios y retorna 404
func noDirectoryListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func New(hub *socket.Hub, config Config) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.ListSpecsHandler(config.Folder))
	r.HandleFunc("/view", handlers.ViewSpecHandler(config.Folder))
	// Servir archivos estáticos sin directory listing
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", noDirectoryListing(http.FileServer(http.Dir("./web/public")))))

	r.HandleFunc("/ws", handlers.WebSocketHandler(hub))

	return &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}
}
