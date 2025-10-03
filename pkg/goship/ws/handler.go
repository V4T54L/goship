package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Default upgrader with sane config
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, tighten this up!
		return true
	},
}

// Handler upgrades the request to WebSocket and calls your handler.
func Handler(fn func(Conn)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusBadRequest)
			return
		}
		conn := &wsConn{Conn: socket}

		// Run user handler
		fn(conn)

		// Ensure connection is closed when done
		conn.Close()
	}
}
