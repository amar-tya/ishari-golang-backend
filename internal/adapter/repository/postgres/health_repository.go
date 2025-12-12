package postgres

import (
	"fmt"

	"ishari-backend/internal/core/port/repository"

	"gorm.io/gorm"
)

type healthRepository struct {
	db *gorm.DB
}

// NewHealthRepository creates a new health repository instance
func NewHealthRepository(db *gorm.DB) repository.HealthRepository {
	return &healthRepository{
		db: db,
	}
}

func (r *healthRepository) CheckDatabase() error {
	if r.db == nil {
		return fmt.Errorf("database not connected")
	}

	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
