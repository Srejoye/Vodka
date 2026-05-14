package main

import (
	"log"

	"github.com/DevanshuTripathi/vodka"
)

// User Type
type User struct {
	Email    string `json:"email" validate:"required,email"` // json and validate struct tags
	Password string `json:"password" validate:"min=8"`
}

func main() {
	app := vodka.DefaultRouter() // Initialize a default router, comes with logger, recovery, error handling middlewares

	app.POST("/create", func(c *vodka.Context) { // POST request on "/create"
		var user User

		err := c.BindJSON(&user) // Binds json fields to struct fields and validates values using validate tags
		if err != nil {
			c.Error(400, err) // c.Error() is core library function to handle errors
			return
		}

		c.JSON(200, vodka.M{ // JSON response returning created user's email
			"message": "user added successfully",
			"email":   user.Email, // user email binded to struct Email field
		})
	})

	if err := app.Run(":8080"); err != nil { // app.Run() starts the server and returns error
		log.Fatalf("Server Didn't Start...")
	}
}

