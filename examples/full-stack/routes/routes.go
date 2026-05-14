package routes

import (
	"full-stack/controllers"

	"github.com/DevanshuTripathi/vodka"
)

// Setup routes
func Setup(r *vodka.Engine) {
	r.GET("/ping", controllers.Pong)

	r.GET("/hello/:name", controllers.Hello)

	r.POST("/create", controllers.CreateNote)

	r.GET("/notes", controllers.GetNotes)
}

