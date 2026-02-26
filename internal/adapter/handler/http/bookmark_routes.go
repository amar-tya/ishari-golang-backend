package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

func RegisterBookmarkRoutes(router fiber.Router, ctrl *controller.BookmarkController, authUC portuc.AuthUseCase) {
	bookmarkGroup := router.Group("/bookmarks")

	// Apply authentication middleware to all bookmark routes
	bookmarkGroup.Use(middleware.AuthMiddleware(authUC))

	// Routes
	bookmarkGroup.Post("/", ctrl.Create)
	bookmarkGroup.Get("/", ctrl.List)
	bookmarkGroup.Get("/user/:userId", ctrl.ListBySpecificUser)
	bookmarkGroup.Get("/:id", ctrl.GetByID)
	bookmarkGroup.Put("/:id", ctrl.Update)
	bookmarkGroup.Delete("/:id", ctrl.Delete)
}
