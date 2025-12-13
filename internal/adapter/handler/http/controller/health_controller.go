package controller

import (
	"ishari-backend/internal/core/usecase"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	healthUseCase usecase.HealthUseCase
}

func NewHealthController(healthUseCase usecase.HealthUseCase) *HealthController {
	return &HealthController{healthUseCase: healthUseCase}
}

func (h *HealthController) Welcome(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (h *HealthController) HealthCheck(c *fiber.Ctx) error {
	status := h.healthUseCase.CheckHealth()
	return c.JSON(status)
}

func (h *HealthController) DatabaseHealthCheck(c *fiber.Ctx) error {
	status := h.healthUseCase.CheckDatabaseHealth()
	if status.Status == "error" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(status)
	}
	return c.JSON(status)
}
