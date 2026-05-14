package main

import (
	"log"
	"time"

	"github.com/DevanshuTripathi/vodka"
)
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

// Custom Logger Middleware
// returns a handler function
func Logger() vodka.HandlerFunc {
	return func(c *vodka.Context) {
		start := time.Now() // Get time before execution

		c.Next() // Process the request through other middlewares

		latency := time.Since(start) // Calculate latency of the request

		log.Printf(
			"["+Blue+"%s"+Reset+"] %s "+Gray+"%v"+Reset,
			c.Request.Method,
			c.Request.URL.Path,
			latency,
		)
	}
}

func main() {
	app := vodka.NewRouter() // Initialize a bare bones router

	app.Use(Logger()) // Use your own custome middleware

	// GET function accepts the endpoint and a handler function
	app.GET("/ping", func(c *vodka.Context) {
		c.String(200, "pong") // Returns a string response with status code 200
	})

	app.GET("/hello/:name", func(c *vodka.Context) {
		name := c.Param("name")
		// Returns JSON response with status code 200
		c.JSON(200, vodka.M{ // vodka.M is shorthand for a Go map
			"message": "Greetings!",
			"name":    name,
		})
	})

	if err := app.Run(":8080"); err != nil { // app.Run() starts the server and returns error
		log.Fatalf("Server Didn't Start...")
	}
}

