package controllers

import (
	"full-stack/models"

	"github.com/DevanshuTripathi/vodka"
)

// Handler Functions for each route

func Pong(c *vodka.Context) {
	c.String(200, "Pong!")
}

func Hello(c *vodka.Context) {
	name := c.Param("name")

	c.String(200, "Hello "+name+"!")
}

func CreateNote(c *vodka.Context) {
	var note models.Note

	err := c.BindJSON(&note)
	if err != nil {
		c.Error(400, err)
		return
	}

	models.Notes = append(models.Notes, note)

	c.JSON(200, vodka.M{
		"message": "successfully created note",
	})
}

func GetNotes(c *vodka.Context) {
	c.JSON(200, models.Notes)
}

