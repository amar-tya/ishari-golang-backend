package entity

import (
	"time"

	"gorm.io/gorm"
)

type Verse struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	ChapterID       uint           `json:"chapter_id" gorm:"not null"`
	Chapter         *Chapter       `json:"chapter,omitempty" gorm:"foreignKey:ChapterID"`
	VerseNumber     uint           `json:"verse_number" gorm:"not null"`
	ArabicText      string         `json:"arabic_text" gorm:"type:text;not null"`
	Transliteration *string        `json:"transliteration,omitempty" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Verse) TableName() string { return "verses" }
