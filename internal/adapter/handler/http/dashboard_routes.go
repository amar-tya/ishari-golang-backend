package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

// RegisterDashboardRoutes registers all dashboard-related HTTP routes
func RegisterDashboardRoutes(router fiber.Router, ctrl *controller.DashboardController, authUC portuc.AuthUseCase) {
	dashboard := router.Group("/dashboard")

	// Protected routes (require JWT token AND super_admin or admin_content roles)
	protected := dashboard.Group("", middleware.AuthMiddleware(authUC))
	protected.Use(middleware.RequireRoles("super_admin", "admin_content"))

	protected.Get("/stats", ctrl.GetStats)
}
