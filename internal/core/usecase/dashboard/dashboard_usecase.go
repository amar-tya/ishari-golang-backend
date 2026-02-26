package dashboard

import (
	"context"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
	"ishari-backend/internal/core/port/usecase"
)

// dashboardUseCase implements usecase.DashboardUseCase
type dashboardUseCase struct {
	repo repository.DashboardRepository
}

// NewDashboardUseCase creates a new dashboard use case instance
func NewDashboardUseCase(repo repository.DashboardRepository) usecase.DashboardUseCase {
	return &dashboardUseCase{repo: repo}
}

func (uc *dashboardUseCase) GetStats(ctx context.Context) (*entity.DashboardStats, error) {
	return uc.repo.GetStats(ctx)
}
