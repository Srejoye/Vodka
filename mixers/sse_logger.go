package mixers

import (
	"log"
	"time"

	"github.com/DevanshuTripathi/vodka"
)

// SSELogger wraps an SSEHandlerFunc and logs connect, disconnect, and duration.
//
// Usage:
//
//	app.SSE("/events", mixers.SSELogger(func(c *vodka.SSEContext) {
//	    // your handler
//	}))
func SSELogger(handler vodka.SSEHandlerFunc) vodka.SSEHandlerFunc {
	return func(c *vodka.SSEContext) {
		start := time.Now()
		log.Printf(vodka.Green+"SSE Connect    %s %s"+vodka.Reset, c.IP(), c.Request.URL.Path)

		handler(c)

		log.Printf(vodka.Blue+"SSE Disconnect %s %s %s"+vodka.Reset, c.IP(), c.Request.URL.Path, time.Since(start))
	}
}
