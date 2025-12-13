package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes registers all user-related HTTP routes
func RegisterUserRoutes(router fiber.Router, ctrl *controller.UserController, authUC portuc.AuthUseCase) {
	users := router.Group("/users")

	// Public routes (no auth required)
	users.Post("/register", ctrl.Register)

	// Protected routes (require JWT token)
	protected := users.Group("", middleware.AuthMiddleware(authUC))
	protected.Get("/", ctrl.ListUsers)
	protected.Get("/:id", ctrl.GetUserByID)
	protected.Put("/:id", ctrl.UpdateUser)
	protected.Delete("/:id", ctrl.DeleteUser)
}
