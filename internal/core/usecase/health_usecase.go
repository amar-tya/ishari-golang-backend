package usecase

import (
	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
)

// HealthUseCase defines the interface for health check use cases
type HealthUseCase interface {
	CheckHealth() *entity.HealthStatus
	CheckDatabaseHealth() *entity.HealthStatus
}

type healthUseCase struct {
	healthRepo repository.HealthRepository
}

// NewHealthUseCase creates a new health use case instance
func NewHealthUseCase(repo repository.HealthRepository) HealthUseCase {
	return &healthUseCase{
		healthRepo: repo,
	}
}

func (uc *healthUseCase) CheckHealth() *entity.HealthStatus {
	return &entity.HealthStatus{
		Status: "ok",
	}
}

func (uc *healthUseCase) CheckDatabaseHealth() *entity.HealthStatus {
	err := uc.healthRepo.CheckDatabase()
	if err != nil {
		return &entity.HealthStatus{
			Status:   "error",
			Database: err.Error(),
		}
	}

	return &entity.HealthStatus{
		Status:   "ok",
		Database: "connected",
	}
}
