package controller

import (
	"ishari-backend/internal/adapter/handler/http/response"
	"ishari-backend/internal/core/port/usecase"
	"ishari-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// DashboardController handles dashboard-related HTTP requests
type DashboardController struct {
	dashboardUseCase usecase.DashboardUseCase
	log              logger.Logger
}

// NewDashboardController creates a new dashboard controller instance
func NewDashboardController(uc usecase.DashboardUseCase, l logger.Logger) *DashboardController {
	return &DashboardController{
		dashboardUseCase: uc,
		log:              l,
	}
}

// GetStats returns the aggregated statistics for the dashboard
// @Summary Get dashboard statistics
// @Description Returns the total users, chapters, verses, hadis, and verse media
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dashboard/stats [get]
func (h *DashboardController) GetStats(c *fiber.Ctx) error {
	stats, err := h.dashboardUseCase.GetStats(c.UserContext())
	if err != nil {
		h.log.Error("Failed to build dashboard stats: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to retrieve dashboard statistics",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response.MapDashboardStatsResponse(stats),
	})
}
