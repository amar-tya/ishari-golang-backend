package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

func RegisterChapterRoutes(router fiber.Router, ctrl *controller.ChapterController, authUC portuc.AuthUseCase) {
	chapter := router.Group("/chapters")

	// Public routes (no auth required)
	chapter.Get("/", ctrl.List)
	chapter.Get("/:id", ctrl.GetByID)
	chapter.Get("/book/:bookId", ctrl.GetByBookID)

	// Protected routes (require JWT token)
	protected := chapter.Group("", middleware.AuthMiddleware(authUC))
	protected.Post("/", ctrl.Create)
	protected.Put("/:id", ctrl.Update)
	protected.Delete("/:id", ctrl.Delete)
	protected.Post("/bulk-delete", ctrl.BulkDelete)
}
