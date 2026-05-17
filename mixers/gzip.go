package mixers

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"

	"github.com/DevanshuTripathi/vodka"
)

// gzipResponseWriter wraps http.ResponseWriter to write compressed data.
type gzipResponseWriter struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.gz.Write(b)
}

// pool reuses gzip writers to avoid allocations on every request.
var gzipPool = sync.Pool{
	New: func() any {
		gz, _ := gzip.NewWriterLevel(nil, gzip.DefaultCompression)
		return gz
	},
}

func newGzipMiddleware(level int) vodka.HandlerFunc {
	return func(c *vodka.Context) {
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz := gzipPool.Get().(*gzip.Writer)
		gz.Reset(c.Writer)

		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Del("Content-Length")

		grw := &gzipResponseWriter{ResponseWriter: c.Writer, gz: gz}
		original := c.Writer
		c.Writer = grw

		defer func() {
			gz.Close()
			gzipPool.Put(gz)
			c.Writer = original
		}()

		c.Next()
	}
}

// Gzip returns a middleware that compresses responses using gzip at default level.
// Use at group or route level so c.Next() runs the route handler inside the middleware.
//
//	api := app.Group("/api", mixers.Gzip())
func Gzip() vodka.HandlerFunc {
	return newGzipMiddleware(gzip.DefaultCompression)
}

// GzipWithLevel returns a Gzip middleware with a custom compression level.
//
//	api := app.Group("/api", mixers.GzipWithLevel(gzip.BestSpeed))
func GzipWithLevel(level int) vodka.HandlerFunc {
	return newGzipMiddleware(level)
}
