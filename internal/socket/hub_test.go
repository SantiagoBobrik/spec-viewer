package socket

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SantiagoBobrik/spec-viewer/pkg/logger"
	"github.com/gorilla/websocket"
)

func init() {
	// Silence logger output during tests.
	logger.SetOutput(io.Discard)
}

// upgrader is used by the test server to upgrade HTTP connections to WebSocket.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// newTestServer creates an httptest server that upgrades connections and adds
// them to the provided Hub. It returns the server and a cleanup function.
func newTestServer(t *testing.T, hub *Hub) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("upgrade error: %v", err)
			return
		}
		hub.Add(conn)
	}))
}

// dial connects to the test server and returns the client-side WebSocket connection.
func dial(t *testing.T, server *httptest.Server) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to dial WebSocket: %v", err)
	}
	return conn
}

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub returned nil")
	}
	if hub.clients == nil {
		t.Fatal("clients map should be initialized")
	}
	if len(hub.clients) != 0 {
		t.Errorf("expected 0 clients, got %d", len(hub.clients))
	}
}

func TestEvents(t *testing.T) {
	if Events.Reload != "reload" {
		t.Errorf("expected Events.Reload to be 'reload', got %q", Events.Reload)
	}
}

func TestHub_AddClient(t *testing.T) {
	hub := NewHub()
	server := newTestServer(t, hub)
	defer server.Close()

	conn := dial(t, server)
	defer func() { _ = conn.Close() }()

	// Give the server handler time to execute Add.
	time.Sleep(50 * time.Millisecond)

	hub.mu.Lock()
	count := len(hub.clients)
	hub.mu.Unlock()

	if count != 1 {
		t.Errorf("expected 1 client, got %d", count)
	}
}

func TestHub_AddMultipleClients(t *testing.T) {
	hub := NewHub()
	server := newTestServer(t, hub)
	defer server.Close()

	conns := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		conns[i] = dial(t, server)
		defer func(c *websocket.Conn) { _ = c.Close() }(conns[i])
	}

	time.Sleep(50 * time.Millisecond)

	hub.mu.Lock()
	count := len(hub.clients)
	hub.mu.Unlock()

	if count != 3 {
		t.Errorf("expected 3 clients, got %d", count)
	}
}

func TestHub_RemoveClient(t *testing.T) {
	hub := NewHub()
	server := newTestServer(t, hub)
	defer server.Close()

	conn := dial(t, server)
	time.Sleep(50 * time.Millisecond)

	// Verify client was added.
	hub.mu.Lock()
	if len(hub.clients) != 1 {
		t.Fatalf("expected 1 client after add, got %d", len(hub.clients))
	}
	// Get the server-side conn from the map.
	var serverConn *websocket.Conn
	for c := range hub.clients {
		serverConn = c
	}
	hub.mu.Unlock()

	// Remove the server-side connection.
	hub.Remove(serverConn)

	hub.mu.Lock()
	count := len(hub.clients)
	hub.mu.Unlock()

	if count != 0 {
		t.Errorf("expected 0 clients after remove, got %d", count)
	}

	// Closing client conn is safe even if server already closed it.
	conn.Close()
}

func TestHub_RemoveNonExistentClient(t *testing.T) {
	hub := NewHub()
	server := newTestServer(t, hub)
	defer server.Close()

	// Create a connection but don't add it to the hub.
	conn := dial(t, server)
	time.Sleep(50 * time.Millisecond)

	// The connection was added by the test server handler, so grab it.
	hub.mu.Lock()
	var serverConn *websocket.Conn
	for c := range hub.clients {
		serverConn = c
	}
	hub.mu.Unlock()

	// Remove it.
	hub.Remove(serverConn)

	// Remove again -- should not panic.
	hub.Remove(serverConn)

	hub.mu.Lock()
	count := len(hub.clients)
	hub.mu.Unlock()

	if count != 0 {
		t.Errorf("expected 0 clients, got %d", count)
	}

	conn.Close()
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	server := newTestServer(t, hub)
	defer server.Close()

	// Connect two clients.
	conn1 := dial(t, server)
	defer func() { _ = conn1.Close() }()
	conn2 := dial(t, server)
	defer func() { _ = conn2.Close() }()

	time.Sleep(50 * time.Millisecond)

	hub.mu.Lock()
	if len(hub.clients) != 2 {
		t.Fatalf("expected 2 clients, got %d", len(hub.clients))
	}
	hub.mu.Unlock()

	// Broadcast a message.
	hub.Broadcast("reload")

	// Both clients should receive the message.
	for i, conn := range []*websocket.Conn{conn1, conn2} {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("client %d: failed to read message: %v", i, err)
		}
		if string(msg) != "reload" {
			t.Errorf("client %d: expected 'reload', got %q", i, string(msg))
		}
	}
}

func TestHub_BroadcastToNoClients(t *testing.T) {
	hub := NewHub()
	// Should not panic with no clients.
	hub.Broadcast("reload")
}

func TestHub_BroadcastRemovesDisconnectedClients(t *testing.T) {
	hub := NewHub()
	server := newTestServer(t, hub)
	defer server.Close()

	conn := dial(t, server)
	time.Sleep(50 * time.Millisecond)

	// Get the server-side connection and close it directly to simulate a broken pipe.
	// This guarantees that WriteMessage will fail on the next Broadcast.
	hub.mu.Lock()
	var serverConn *websocket.Conn
	for c := range hub.clients {
		serverConn = c
	}
	hub.mu.Unlock()

	// Close the underlying network connection to force a write failure.
	_ = serverConn.UnderlyingConn().Close()

	hub.Broadcast("reload")

	// The hub should have removed the dead client.
	hub.mu.Lock()
	count := len(hub.clients)
	hub.mu.Unlock()

	if count != 0 {
		t.Errorf("expected 0 clients after broadcast to dead conn, got %d", count)
	}

	conn.Close()
}
