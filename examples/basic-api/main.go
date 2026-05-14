package main

import (
	"log"

	"github.com/DevanshuTripathi/vodka"
)

func main() {
	app := vodka.DefaultRouter() // Initialize a default router, comes with logger, recovery, error handling middlewares

	// GET function accepts the endpoint and a handler function
	app.GET("/ping", func(c *vodka.Context) {
		c.String(200, "pong") // Returns a string response with status code 200
	})

	app.GET("/hello/:name", func(c *vodka.Context) {
		name := c.Param("name") // gets url param values
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

