package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"

	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes registers all user-related HTTP routes
func RegisterUserRoutes(router fiber.Router, ctrl *controller.UserController) {
	users := router.Group("/users")

	// Public routes
	users.Post("/register", ctrl.Register)
	users.Post("/login", ctrl.Login)

	// Protected routes (middleware can be added later)
	users.Get("/", ctrl.ListUsers)
	users.Get("/:id", ctrl.GetUserByID)
	users.Put("/:id", ctrl.UpdateUser)
	users.Delete("/:id", ctrl.DeleteUser)
}
