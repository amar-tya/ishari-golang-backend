package entity

import (
	"time"
)

// DashboardStats represents the aggregated statistics for the dashboard
type DashboardStats struct {
	TotalUsers      int       `json:"total_users"`
	TotalHadis      int       `json:"total_hadis"`
	TotalChapters   int       `json:"total_chapters"`
	TotalVerses     int       `json:"total_verses"`
	TotalVerseMedia int       `json:"total_verse_media"`
	CalculatedAt    time.Time `json:"calculated_at"`
}
