package main

import (
	"log"

	"github.com/DevanshuTripathi/vodka"
	"github.com/DevanshuTripathi/vodka/mixers"
	"github.com/gorilla/websocket"
)

func main() {
	app := vodka.DefaultRouter()
	app.Use(vodka.AllowCORS([]string{"*"}))
	app.AllowWSOrigins([]string{"http://localhost:8080", "https://userapp.com"})

	// Basic echo — reflects every message back to the sender.
	app.WS("/ws/echo", mixers.WSLogger(func(c *vodka.WSContext) {
		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				// CloseError means the client disconnected cleanly.
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Printf("read error: %v", err)
				}
				return
			}
			if err := c.WriteMessage(msgType, msg); err != nil {
				log.Printf("write error: %v", err)
				return
			}
		}
	}))

	// Room-based broadcast — clients in the same :room receive each other's messages.
	hub := newHub()
	go hub.run()

	app.WS("/ws/room/:room", mixers.WSLogger(func(c *vodka.WSContext) {
		room := c.Param("room")
		client := &client{hub: hub, conn: c.Conn, room: room, send: make(chan []byte, 256)}
		hub.register <- client

		go client.writePump()
		client.readPump()
	}))

	// JSON ping-pong — demonstrates WriteJSON / ReadJSON.
	app.WS("/ws/json", func(c *vodka.WSContext) {
		type Ping struct {
			Message string `json:"message"`
		}
		type Pong struct {
			Reply string `json:"reply"`
		}

		for {
			var p Ping
			if err := c.ReadJSON(&p); err != nil {
				return
			}
			if err := c.WriteJSON(Pong{Reply: "pong: " + p.Message}); err != nil {
				return
			}
		}
	})

	if err := app.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Minimal broadcast hub
// ---------------------------------------------------------------------------

type client struct {
	hub  *hub
	conn *websocket.Conn
	room string
	send chan []byte
}

func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		c.hub.broadcast <- &message{room: c.room, data: msg}
	}
}

func (c *client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

type message struct {
	room string
	data []byte
}

type hub struct {
	clients    map[*client]bool
	broadcast  chan *message
	register   chan *client
	unregister chan *client
}

func newHub() *hub {
	return &hub{
		clients:    make(map[*client]bool),
		broadcast:  make(chan *message),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.clients {
				if c.room != m.room {
					continue
				}
				select {
				case c.send <- m.data:
				default:
					delete(h.clients, c)
					close(c.send)
				}
			}
		}
	}
}
