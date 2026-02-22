package entity

import (
	"time"
)

type Book struct {
	ID            int        `json:"id" gorm:"primaryKey"`
	Title         string     `json:"title"`
	Author        *string    `json:"author,omitempty"`
	Description   *string    `json:"description,omitempty"`
	PublishedYear *int       `json:"published_year,omitempty"`
	CoverImageURL *string    `json:"cover_image_url,omitempty"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `json:"-" gorm:"index"`
}

func (Book) TableName() string { return "books" }
