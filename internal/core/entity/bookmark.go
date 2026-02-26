package entity

import (
	"time"

	"gorm.io/gorm"
)

type Bookmark struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	VerseID   uint           `json:"verse_id" gorm:"not null"`
	Note      *string        `json:"note" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Bookmark) TableName() string {
	return "bookmarks"
}
