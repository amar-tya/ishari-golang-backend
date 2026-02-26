package response

import (
	"time"

	"ishari-backend/internal/core/entity"
)

type DashboardStatsResponse struct {
	TotalUsers      int       `json:"total_users"`
	TotalHadis      int       `json:"total_hadis"`
	TotalChapters   int       `json:"total_chapters"`
	TotalVerses     int       `json:"total_verses"`
	TotalVerseMedia int       `json:"total_verse_media"`
	CalculatedAt    time.Time `json:"calculated_at"`
}

func MapDashboardStatsResponse(stats *entity.DashboardStats) DashboardStatsResponse {
	return DashboardStatsResponse{
		TotalUsers:      stats.TotalUsers,
		TotalHadis:      stats.TotalHadis,
		TotalChapters:   stats.TotalChapters,
		TotalVerses:     stats.TotalVerses,
		TotalVerseMedia: stats.TotalVerseMedia,
		CalculatedAt:    stats.CalculatedAt,
	}
}
