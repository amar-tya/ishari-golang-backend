package usecase

import (
	"context"

	"ishari-backend/internal/core/entity"
)

// DashboardUseCase defines the business logic for the dashboard
type DashboardUseCase interface {
	// GetStats returns aggregated statistics for the entire platform
	GetStats(ctx context.Context) (*entity.DashboardStats, error)
}
