package mixers

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DevanshuTripathi/vodka"
)

func TestGzipCompressesResponse(t *testing.T) {
	app := vodka.NewRouter()
	api := app.Group("/api", Gzip())
	api.GET("/data", func(c *vodka.Context) {
		c.JSON(200, vodka.M{"message": "hello"})
	})

	s := httptest.NewServer(app)
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL+"/api/data", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("expected Content-Encoding: gzip, got %q", resp.Header.Get("Content-Encoding"))
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("gzip reader error: %v", err)
	}
	defer gz.Close()

	body, _ := io.ReadAll(gz)
	if !strings.Contains(string(body), "hello") {
		t.Errorf("decompressed body missing expected content, got: %s", body)
	}
}

func TestGzipSkipsWhenNotAccepted(t *testing.T) {
	app := vodka.NewRouter()
	api := app.Group("/api", Gzip())
	api.GET("/data", func(c *vodka.Context) {
		c.JSON(200, vodka.M{"message": "hello"})
	})

	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL + "/api/data")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if enc := resp.Header.Get("Content-Encoding"); enc == "gzip" {
		t.Error("should not gzip when Accept-Encoding header is absent")
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "hello") {
		t.Errorf("body missing expected content, got: %s", body)
	}
}

func TestGzipWithLevel(t *testing.T) {
	app := vodka.NewRouter()
	api := app.Group("/api", GzipWithLevel(gzip.BestSpeed))
	api.GET("/data", func(c *vodka.Context) {
		c.String(200, "fast compression")
	})

	s := httptest.NewServer(app)
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL+"/api/data", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("expected Content-Encoding: gzip, got %q", resp.Header.Get("Content-Encoding"))
	}

	gz, _ := gzip.NewReader(resp.Body)
	defer gz.Close()
	body, _ := io.ReadAll(gz)
	if string(body) != "fast compression" {
		t.Errorf("got %q, want %q", body, "fast compression")
	}
}

func TestGzipRouteLevel(t *testing.T) {
	app := vodka.NewRouter()
	g := app.Group("", Gzip())
	g.GET("/compressed", func(c *vodka.Context) {
		c.JSON(200, vodka.M{"ok": true})
	})

	s := httptest.NewServer(app)
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL+"/compressed", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("expected Content-Encoding: gzip at route level")
	}
}

func TestGzipPreservesStatusCode(t *testing.T) {
	app := vodka.NewRouter()
	app.Use(vodka.Logger())
	api := app.Group("/api", Gzip())
	api.GET("/notfound", func(c *vodka.Context) {
		c.JSON(http.StatusNotFound, vodka.M{"error": "not found"})
	})

	s := httptest.NewServer(app)
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL+"/api/notfound", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGzipPreservesStatusCodeOnError(t *testing.T) {
	app := vodka.NewRouter()
	app.Use(vodka.Logger())
	api := app.Group("/api", Gzip())
	api.GET("/error", func(c *vodka.Context) {
		c.JSON(http.StatusInternalServerError, vodka.M{"error": "server error"})
	})

	s := httptest.NewServer(app)
	defer s.Close()

	req, _ := http.NewRequest("GET", s.URL+"/api/error", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}