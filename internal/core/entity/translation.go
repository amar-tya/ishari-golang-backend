package entity

import (
	"time"

	"gorm.io/gorm"
)

type Translation struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	VerseID         uint           `json:"verse_id" gorm:"not null"`
	LanguageCode    string         `json:"language_code" gorm:"type:varchar(10);not null"`
	TranslationText string         `json:"translation_text" gorm:"type:text;not null"`
	TranslatorName  *string        `json:"translator_name,omitempty" gorm:"type:varchar(255)"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Translation) TableName() string { return "translations" }
