package dashboard_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/usecase/dashboard"
)

// MockDashboardRepository for testing
type MockDashboardRepository struct {
	GetStatsFunc func(ctx context.Context) (*entity.DashboardStats, error)
}

func (m *MockDashboardRepository) GetStats(ctx context.Context) (*entity.DashboardStats, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx)
	}
	return &entity.DashboardStats{
		TotalUsers:      10,
		TotalHadis:      5,
		TotalChapters:   114,
		TotalVerses:     6236,
		TotalVerseMedia: 100,
		CalculatedAt:    time.Now(),
	}, nil
}

func TestDashboardUseCase_GetStats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedStats := &entity.DashboardStats{
			TotalUsers:      20,
			TotalHadis:      10,
			TotalChapters:   114,
			TotalVerses:     6236,
			TotalVerseMedia: 200,
			CalculatedAt:    time.Now(),
		}

		mockRepo := &MockDashboardRepository{
			GetStatsFunc: func(ctx context.Context) (*entity.DashboardStats, error) {
				return expectedStats, nil
			},
		}

		uc := dashboard.NewDashboardUseCase(mockRepo)
		stats, err := uc.GetStats(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if stats.TotalUsers != expectedStats.TotalUsers {
			t.Errorf("expected TotalUsers %d, got %d", expectedStats.TotalUsers, stats.TotalUsers)
		}
		if stats.TotalHadis != expectedStats.TotalHadis {
			t.Errorf("expected TotalHadis %d, got %d", expectedStats.TotalHadis, stats.TotalHadis)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := &MockDashboardRepository{
			GetStatsFunc: func(ctx context.Context) (*entity.DashboardStats, error) {
				return nil, errors.New("database error")
			},
		}

		uc := dashboard.NewDashboardUseCase(mockRepo)
		_, err := uc.GetStats(context.Background())

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if err.Error() != "database error" {
			t.Errorf("expected error 'database error', got '%v'", err)
		}
	})
}
