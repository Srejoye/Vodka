package main

import (
	"full-stack/routes"
	"log"

	"github.com/DevanshuTripathi/vodka"
)

func main() {
	app := vodka.DefaultRouter() // Creates a Default Router with Logger and Recovery Middleware

	allowedOrigins := []string{"http://localhost:5173"} // Allow vite frontend

	app.Use(vodka.AllowCORS(allowedOrigins)) // Add the core AllowCORS middleware

	routes.Setup(app)

	if err := app.Run(":8080"); err != nil { // app.Run() starts the server and returns error
		log.Fatalf("Server Didn't Start...")
	}
}

