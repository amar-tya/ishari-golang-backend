package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"

	"github.com/gofiber/fiber/v2"
)

func RegisterHealthRoutes(router fiber.Router, h *controller.HealthController) {
	router.Get("/", h.Welcome)
	router.Get("/health", h.HealthCheck)
	router.Get("/health/db", h.DatabaseHealthCheck)
}
