package vodka

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// WSHandlerFunc is the handler signature for WebSocket connections.
type WSHandlerFunc func(*WSContext)

// WSConfig holds configuration for the WebSocket upgrader.
type WSConfig struct {
	ReadBufferSize   int
	WriteBufferSize  int
	HandshakeTimeout time.Duration
	// CheckOrigin returns true to allow the connection. Defaults to allowing all origins.
	CheckOrigin func(r *http.Request) bool
}

// DefaultWSConfig returns a WSConfig with sensible defaults.
// By default only same-origin requests are allowed. Use engine.AllowWSOrigins() to permit others.
func DefaultWSConfig() *WSConfig {
	return &WSConfig{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 10 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			// Allow requests with no Origin header (non-browser clients).
			if origin == "" {
				return true
			}
			// Only allow same-origin by default.
			return origin == "http://"+r.Host || origin == "https://"+r.Host
		},
	}
}

// WSContext wraps a live WebSocket connection with framework utilities.
// URL params, query values, and keys set by HTTP middlewares are all accessible here.
type WSContext struct {
	Conn    *websocket.Conn
	Keys    map[string]any
	Params  httprouter.Params
	Request *http.Request
}

// Set stores a value in the context by key, shared across the connection lifecycle.
func (wc *WSContext) Set(key string, value any) {
	if wc.Keys == nil {
		wc.Keys = make(map[string]any)
	}
	wc.Keys[key] = value
}

// Get retrieves a value previously stored via Set.
func (wc *WSContext) Get(key string) (value any, exists bool) {
	value, exists = wc.Keys[key]
	return
}

// Param returns a URL path parameter by name (e.g. /ws/:room → wc.Param("room")).
func (wc *WSContext) Param(key string) string {
	return wc.Params.ByName(key)
}

// Query returns a URL query parameter by name.
func (wc *WSContext) Query(key string) string {
	return wc.Request.URL.Query().Get(key)
}

// IP returns the remote address of the client.
func (wc *WSContext) IP() string {
	return wc.Request.RemoteAddr
}

// ReadMessage reads the next message from the connection.
// messageType is websocket.TextMessage or websocket.BinaryMessage.
func (wc *WSContext) ReadMessage() (messageType int, p []byte, err error) {
	return wc.Conn.ReadMessage()
}

// WriteMessage writes a message to the connection.
func (wc *WSContext) WriteMessage(messageType int, data []byte) error {
	return wc.Conn.WriteMessage(messageType, data)
}

// WriteJSON encodes v as JSON and sends it as a text message.
func (wc *WSContext) WriteJSON(v any) error {
	return wc.Conn.WriteJSON(v)
}

// ReadJSON reads the next message and decodes it as JSON into v.
func (wc *WSContext) ReadJSON(v any) error {
	return wc.Conn.ReadJSON(v)
}

// Close sends a normal closure message and closes the connection gracefully.
func (wc *WSContext) Close() error {
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	return wc.Conn.WriteMessage(websocket.CloseMessage, msg)
}
