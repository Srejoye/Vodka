package main

import (
	"errors"
	"log"

	"github.com/DevanshuTripathi/vodka"
)

func main() {
	app := vodka.NewRouter() // Initialize a bare bones router

	app.Use(vodka.Logger(), vodka.Recovery())

	// GET function accepts the endpoint and a handler function
	app.GET("/ping", func(c *vodka.Context) {
		c.String(200, "pong") // Returns a string response with status code 200
	})

	// Create a router group with middlewares
	// The middleware will only be added to the routes which come under the router group
	// Another way of creating a group is
	// api := app.Group("/api")
	// and add middleware later with
	// api.Use(vodka.ErrorHandler())
	api := app.Group("/api", vodka.ErrorHandler())

	// GET request on url: "/api/user/:userId"
	api.GET("/user/:userId", func(c *vodka.Context) {
		userId := c.Param("userId") // gets url param values
		// Returns JSON response with status code 200
		c.JSON(200, vodka.M{ // vodka.M is shorthand for a Go map
			"userId": userId,
			"name":   "generic username",
		})
	})

	// GET request for error example
	api.GET("/error", func(c *vodka.Context) {
		c.Error(500, errors.New("example error")) // c.Error for handling error with ErrorHandler middleware
	})

	if err := app.Run(":8080"); err != nil { // app.Run() starts the server and returns error
		log.Fatalf("Server Didn't Start...")
	}
}

