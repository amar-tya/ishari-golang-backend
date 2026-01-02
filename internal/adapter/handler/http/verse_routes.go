package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

func RegisterVerseRoutes(router fiber.Router, ctrl *controller.VerseController, authUC portuc.AuthUseCase) {
	verse := router.Group("/verses")

	// Public routes (no auth required)
	verse.Get("/", ctrl.List)
	verse.Get("/:id", ctrl.GetByID)

	// Protected routes (require JWT token)
	protected := verse.Group("", middleware.AuthMiddleware(authUC))
	protected.Post("/", ctrl.Create)
	protected.Put("/:id", ctrl.Update)
	protected.Delete("/:id", ctrl.Delete)
}
