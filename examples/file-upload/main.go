package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/DevanshuTripathi/vodka"
)

const maxUploadSize = 10 << 20 // 10 MB

func main() {
	app := vodka.DefaultRouter()

	// Health check endpoint
	app.GET("/health", func(c *vodka.Context) {
		c.JSON(200, vodka.M{
			"status":  "ok",
			"service": "file-upload-example",
		})
	})

	// POST /upload — accepts a multipart file upload, saves it to ./uploads,
	// and returns non-sensitive metadata about the saved file.
	app.POST("/upload", func(c *vodka.Context) {
		// Enforce 10 MB total request body cap before any multipart parsing.
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)
		if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
			log.Printf("[upload] multipart parse failed (possible size-limit exceeded): %v", err)
			c.Error(400, fmt.Errorf("request body exceeds 10 MB limit or is malformed"))
			return
		}

		// Ensure the uploads directory exists
		uploadsDir := "./uploads"
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Printf("[upload] failed to create uploads directory: %v", err)
			c.Error(500, fmt.Errorf("could not prepare upload directory"))
			return
		}

		// Retrieve the file from the already-parsed multipart form (field name: "file")
		fileHeader, err := c.FormFile("file")
		if err != nil {
			log.Printf("[upload] failed to retrieve file from form: %v", err)
			c.Error(400, fmt.Errorf("missing or invalid file field — use form key \"file\""))
			return
		}

		// Build a unique destination path to avoid name collisions
		timestamp := time.Now().UnixNano()
		safeName := filepath.Base(fileHeader.Filename)
		dstPath := filepath.Join(uploadsDir, fmt.Sprintf("%d_%s", timestamp, safeName))

		// Stream the uploaded file to disk using the built-in Vodka helper
		if err := c.SaveUploadedFile(fileHeader, dstPath); err != nil {
			log.Printf("[upload] failed to save file %q: %v", safeName, err)
			c.Error(500, fmt.Errorf("failed to save uploaded file"))
			return
		}

		log.Printf("[upload] saved %q (%d bytes) → %s", safeName, fileHeader.Size, dstPath)

		// Return only non-sensitive metadata; internal path is intentionally omitted.
		c.JSON(200, vodka.M{
			"message":  "file uploaded successfully",
			"filename": safeName,
			"size":     fileHeader.Size,
		})
	})

	log.Println("File-upload example starting on :8080")
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
