package mixers

import (
	"strings"

	"github.com/DevanshuTripathi/vodka"
)

// Token Validator type recieves token and returns the data linked to it and if the token is valid
type TokenValidator func(token string) (data any, isValid bool)

// Bearer Auth for getting the bearer token from the request header
// Then validate the token using the validator function
// Set value in the context for later use from a middleware
func BearerAuth(ctxKey string, validator TokenValidator) vodka.HandlerFunc {
	return func(c *vodka.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, vodka.M{"error": "Unauthorized: Missing or malformed token"})
			c.Abort()
			return
		}

		providedToken := strings.TrimPrefix(authHeader, "Bearer ")

		data, isValid := validator(providedToken)

		if !isValid {
			c.JSON(401, vodka.M{"error": "Unauthorized: Invalid token"})
			c.Abort()
			return
		}

		c.Set(ctxKey, data)

		c.Next()
	}
}
