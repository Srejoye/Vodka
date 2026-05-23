package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DevanshuTripathi/vodka"
	"github.com/DevanshuTripathi/vodka/mixers"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret is the HMAC signing key. Loaded from JWT_SECRET env var at startup.
// Falls back to a labelled insecure value for local development only.
var jwtSecret = func() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	log.Println("[warn] JWT_SECRET not set — using INSECURE hardcoded fallback (development only)")
	return "super-secret-vodka-key-INSECURE-DO-NOT-USE-IN-PROD"
}()

// loginRequest represents the expected JSON body for POST /login.
type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func main() {
	app := vodka.DefaultRouter()

	// ──────────────────────────────────────────────
	// Public routes
	// ──────────────────────────────────────────────

	// Health check
	app.GET("/health", func(c *vodka.Context) {
		c.JSON(200, vodka.M{
			"status":  "ok",
			"service": "jwt-auth-example",
		})
	})

	// POST /login — authenticates a user and returns a signed JWT.
	// For demonstration purposes any username with password "vodka" is accepted.
	app.POST("/login", func(c *vodka.Context) {
		var req loginRequest
		if err := c.BindJSON(&req); err != nil {
			log.Printf("[login] JSON bind/validation error: %v", err)
			c.Error(400, fmt.Errorf("invalid request body: validation failed"))
			return
		}

		// Simple credential check (swap with a real DB lookup in production)
		if req.Password != "vodka" {
			c.Error(401, fmt.Errorf("invalid credentials"))
			return
		}

		// Build the JWT payload and sign the token using Vodka's mixers helper
		token, err := mixers.GenerateJWT(jwtSecret, map[string]any{
			"sub":      req.Username,
			"username": req.Username,
			"role":     "user",
		}, 24*time.Hour)
		if err != nil {
			log.Printf("[login] token generation failed: %v", err)
			c.Error(500, fmt.Errorf("could not generate token"))
			return
		}

		log.Printf("[login] issued token for user %q", req.Username)
		c.JSON(200, vodka.M{
			"message":    "login successful",
			"token":      token,
			"token_type": "Bearer",
			"expires_in": "24h",
		})
	})

	// ──────────────────────────────────────────────
	// Protected route group — /api/secure/*
	// ──────────────────────────────────────────────

	// Create a JWT token validator using Vodka's built-in helper.
	// BearerAuth extracts the "Authorization: Bearer <token>" header,
	// validates the token, and stores the decoded claims in the context
	// under the key "claims".
	jwtMiddleware := mixers.BearerAuth("claims", mixers.JWTValidator(jwtSecret))

	// Create a RouterGroup with the JWT middleware applied to every sub-route.
	secure := app.Group("/api/secure", jwtMiddleware)

	// GET /api/secure/profile — returns the authenticated user's profile
	// extracted from JWT claims.
	secure.GET("/profile", func(c *vodka.Context) {
		raw, exists := c.Get("claims")
		if !exists {
			c.Error(401, fmt.Errorf("claims not found in context"))
			return
		}
		claims, ok := raw.(jwt.MapClaims)
		if !ok {
			c.Error(500, fmt.Errorf("unexpected claims format"))
			return
		}
		c.JSON(200, vodka.M{
			"message":  "welcome to your profile",
			"username": claims["username"],
			"role":     claims["role"],
			"issued":   claims["iat"],
			"expires":  claims["exp"],
		})
	})

	// GET /api/secure/dashboard — a second protected route to demonstrate
	// the group middleware protecting multiple endpoints.
	secure.GET("/dashboard", func(c *vodka.Context) {
		raw, exists := c.Get("claims")
		if !exists {
			c.Error(401, fmt.Errorf("claims not found in context"))
			return
		}
		claims, ok := raw.(jwt.MapClaims)
		if !ok {
			c.Error(500, fmt.Errorf("unexpected claims format"))
			return
		}
		c.JSON(200, vodka.M{
			"message": "dashboard data",
			"user":    claims["username"],
			"stats": vodka.M{
				"projects":   12,
				"tasks_done": 47,
				"uptime":     "99.9%",
			},
		})
	})

	log.Println("JWT-auth example starting on :8080")
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
