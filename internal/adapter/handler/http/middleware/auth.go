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

		// Propagate claims to UserContext for usecase layer
		ctx := portuc.NewContextWithUser(c.UserContext(), claims)
		c.SetUserContext(ctx)

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

// RequireRoles restricts access to specific roles
func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "unauthorized",
			})
		}

		// Allow super_admin by default (optional, based on your business logic)
		if user.Role == "super_admin" {
			return c.Next()
		}

		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "forbidden: insufficient permissions",
		})
	}
}
