package server

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/SantiagoBobrik/spec-viewer/internal/handlers"
	"github.com/SantiagoBobrik/spec-viewer/internal/socket"
	"github.com/SantiagoBobrik/spec-viewer/web"

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
			handlers.NotFoundHandler().ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func New(hub *socket.Hub, config Config) *http.Server {
	r := mux.NewRouter()

	r.NotFoundHandler = handlers.NotFoundHandler()

	r.HandleFunc("/", handlers.HomeHandler())
	r.HandleFunc("/view", handlers.ViewSpecHandler(config.Folder))

	publicFS, err := fs.Sub(web.Files, "public")
	if err != nil {
		log.Fatalf("Error creating public filesystem: %v", err)
	}

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", noDirectoryListing(http.FileServer(http.FS(publicFS)))))

	r.HandleFunc("/ws", handlers.WebSocketHandler(hub))

	return &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}
}
