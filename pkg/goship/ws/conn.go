package ws

import (
	"github.com/gorilla/websocket"
)

// Conn is a simple wrapper around the gorilla websocket.Conn
type Conn interface {
	ReadJSON(v any) error
	WriteJSON(v any) error
	Close() error
}

type wsConn struct {
	*websocket.Conn
}

func (c *wsConn) ReadJSON(v any) error {
	return c.Conn.ReadJSON(v)
}

func (c *wsConn) WriteJSON(v any) error {
	return c.Conn.WriteJSON(v)
}

func (c *wsConn) Close() error {
	return c.Conn.Close()
}
