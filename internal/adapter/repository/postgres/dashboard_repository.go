package postgres

import (
	"context"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"

	"gorm.io/gorm"
)

// dashboardRepository implements repository.DashboardRepository using GORM
type dashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository(db *gorm.DB) repository.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(ctx context.Context) (*entity.DashboardStats, error) {
	var stats entity.DashboardStats

	query := `
		SELECT 
			(SELECT COUNT(*) FROM users WHERE deleted_at IS NULL) as total_users,
			(SELECT COUNT(*) FROM hadi WHERE deleted_at IS NULL) as total_hadis,
			(SELECT COUNT(*) FROM chapters WHERE deleted_at IS NULL) as total_chapters,
			(SELECT COUNT(*) FROM verses WHERE deleted_at IS NULL) as total_verses,
			(SELECT COUNT(*) FROM verse_media WHERE deleted_at IS NULL) as total_verse_media
	`
	// The table name dual isn't needed in Postgres for a SELECT without FROM,
	// but GORM Raw executes it fine.
	err := r.db.WithContext(ctx).Raw(query).Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	stats.CalculatedAt = time.Now()
	return &stats, nil
}
