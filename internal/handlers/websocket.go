package handlers

import (
	"net/http"

	"github.com/SantiagoBobrik/spec-viewer/internal/socket"
	"github.com/SantiagoBobrik/spec-viewer/pkg/logger"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local development
	},
}

func WebSocketHandler(hub *socket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("Failed to upgrade to websocket", "error", err)
			return
		}
		hub.Add(conn)

		// Start a goroutine to keep the connection open and detect disconnects
		go func() {
			defer hub.Remove(conn)
			for {
				// Although we don't expect messages from the client, reading is necessary
				// to process control frames (ping/pong) and detect when the client closes the connection.
				_, _, err := conn.ReadMessage()
				if err != nil {
					// This is expected when the client disconnects
					break
				}
			}
		}()
	}
}
