package server

import (
	"net/http"
	"spec-viewer/internal/handlers"
	"spec-viewer/internal/socket"

	"github.com/gorilla/mux"
)

type Config struct {
	Port   string
	Folder string
}

func New(hub *socket.Hub, config Config) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler)
	r.PathPrefix("/specs/").Handler(http.StripPrefix("/specs/", http.FileServer(http.Dir(config.Folder))))
	r.HandleFunc("/ws", handlers.WebSocketHandler(hub))

	return &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}
}
