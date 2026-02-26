package repository

import (
	"context"

	"ishari-backend/internal/core/entity"
)

// DashboardRepository provides data access for dashboard statistics
type DashboardRepository interface {
	// GetStats retrieves global application statistics efficiently
	GetStats(ctx context.Context) (*entity.DashboardStats, error)
}
