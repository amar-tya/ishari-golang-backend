package entity

import (
	"time"

	"gorm.io/gorm"
)

type Hadi struct {
	ID          int            `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	ImageURL    *string        `json:"image_url,omitempty"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Hadi) TableName() string { return "hadi" }
