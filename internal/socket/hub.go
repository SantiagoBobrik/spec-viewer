package socket

import (
	"sync"

	"github.com/SantiagoBobrik/spec-viewer/pkg/logger"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

var Events = struct {
	Reload string
}{
	Reload: "reload",
}

func (h *Hub) Add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = true
	logger.Info("Client connected", "total_clients", len(h.clients))
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[conn]; ok {
		delete(h.clients, conn)
		conn.Close()
		logger.Info("Client disconnected", "total_clients", len(h.clients))
	}
}

func (h *Hub) Broadcast(message string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	payload := []byte(message)
	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			logger.Error("Failed to write to websocket, removing client", "error", err)
			delete(h.clients, conn)
			conn.Close()
		}
	}
}
