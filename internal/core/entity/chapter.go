package entity

import (
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	BookID        uint           `json:"book_id" gorm:"not null"`
	ChapterNumber uint           `json:"chapter_number" gorm:"not null"`
	Title         string         `json:"title" gorm:"not null"`
	Category      string         `json:"category" gorm:"not null"`
	Description   *string        `json:"description,omitempty"`
	TotalVerses   uint           `json:"total_verses" gorm:"default:0"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Chapter) TableName() string { return "chapters" }
