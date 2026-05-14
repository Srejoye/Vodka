package models

// Define Note model
type Note struct {
	Content string `json:"content"`
}

// Create in memory store
var Notes = make([]Note, 0)

