package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

// RegisterAuthRoutes registers all authentication-related HTTP routes
func RegisterAuthRoutes(router fiber.Router, ctrl *controller.AuthController, authUC portuc.AuthUseCase) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/login", ctrl.Login)
	auth.Post("/refresh", ctrl.RefreshToken)

	// Protected routes
	auth.Post("/logout", middleware.AuthMiddleware(authUC), ctrl.Logout)
	auth.Get("/me", middleware.AuthMiddleware(authUC), ctrl.GetMe)
}
