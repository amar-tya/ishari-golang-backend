package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

func RegisterHadiRoutes(api fiber.Router, ctrl *controller.HadiController, authUC portuc.AuthUseCase) {
	hadis := api.Group("/hadis")

	// Public routes
	hadis.Get("/", ctrl.List)
	hadis.Get("/:id", ctrl.GetByID)

	// Protected block
	adminBlock := hadis.Group("")
	adminBlock.Use(middleware.AuthMiddleware(authUC))
	adminBlock.Use(middleware.RequireRoles("super_admin", "admin_content"))

	adminBlock.Post("/", ctrl.Create)
	adminBlock.Put("/:id", ctrl.Update)
	adminBlock.Delete("/:id", ctrl.Delete)
}
