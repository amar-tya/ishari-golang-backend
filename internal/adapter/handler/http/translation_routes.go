package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

func RegisterTranslationRoutes(router fiber.Router, ctrl *controller.TranslationController, authUC portuc.AuthUseCase) {
	translation := router.Group("/translations")

	// Public routes (no auth required)
	translation.Get("/", ctrl.List)
	translation.Get("/:id", ctrl.GetByID)
	translation.Get("/verse/:verse_id", ctrl.GetByVerseID)

	// Protected routes (require JWT token)
	protected := translation.Group("", middleware.AuthMiddleware(authUC))
	protected.Post("/", ctrl.Create)
	protected.Put("/:id", ctrl.Update)
	protected.Delete("/:id", ctrl.Delete)
}
