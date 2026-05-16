package main

import (
	"log"
	"time"

	"github.com/DevanshuTripathi/vodka"
	"github.com/DevanshuTripathi/vodka/mixers"
)

func main() {
	app := vodka.DefaultRouter()
	app.Use(vodka.AllowCORS([]string{"*"}))

	// Live clock — pushes the current time every second until the client disconnects.
	app.SSE("/events/clock", mixers.SSELogger(func(c *vodka.SSEContext) {
		for {
			select {
			case <-c.Done():
				return
			default:
				err := c.Send("tick", vodka.M{
					"time": time.Now().Format(time.RFC3339),
				})
				if err != nil {
					return
				}
				time.Sleep(time.Second)
			}
		}
	}))

	// Counter — pushes an incrementing number every 500ms.
	app.SSE("/events/counter", func(c *vodka.SSEContext) {
		i := 0
		for {
			select {
			case <-c.Done():
				return
			default:
				if err := c.Send("count", vodka.M{"value": i}); err != nil {
					return
				}
				i++
				time.Sleep(500 * time.Millisecond)
			}
		}
	})

	// Parameterised channel — /events/feed/:topic
	app.SSE("/events/feed/:topic", func(c *vodka.SSEContext) {
		topic := c.Param("topic")
		for {
			select {
			case <-c.Done():
				return
			default:
				if err := c.Send("message", vodka.M{
					"topic":   topic,
					"payload": "update for " + topic,
					"ts":      time.Now().Unix(),
				}); err != nil {
					return
				}
				time.Sleep(2 * time.Second)
			}
		}
	})

	if err := app.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
