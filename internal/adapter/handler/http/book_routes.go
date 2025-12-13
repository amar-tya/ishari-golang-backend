package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

// RegisterBookRoutes registers all book-related HTTP routes
func RegisterBookRoutes(router fiber.Router, ctrl *controller.BookController, authUC portuc.AuthUseCase) {
	books := router.Group("/books")

	// Public routes (read-only)
	books.Get("/", ctrl.ListBooks)
	books.Get("/:id", ctrl.GetBookById)

	// Protected routes (require JWT token for mutations)
	protected := books.Group("", middleware.AuthMiddleware(authUC))
	protected.Post("/", ctrl.CreateBook)
	protected.Put("/:id", ctrl.EditBook)
	protected.Delete("/:id", ctrl.DeleteBook)
}
