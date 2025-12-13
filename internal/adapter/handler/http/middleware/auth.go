package middleware

import (
	"strings"

	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(authUC portuc.AuthUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "missing authorization header",
			})
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := authUC.ValidateToken(c.UserContext(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		// Store claims in context locals for handlers to access
		c.Locals("user", claims)
		c.Locals("token", token)

		return c.Next()
	}
}

// GetUserFromContext retrieves user claims from fiber context
func GetUserFromContext(c *fiber.Ctx) *portuc.TokenClaims {
	claims, ok := c.Locals("user").(*portuc.TokenClaims)
	if !ok {
		return nil
	}
	return claims
}

// GetTokenFromContext retrieves the raw token from fiber context
func GetTokenFromContext(c *fiber.Ctx) string {
	token, ok := c.Locals("token").(string)
	if !ok {
		return ""
	}
	return token
}
