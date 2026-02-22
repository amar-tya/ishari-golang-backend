package usecase

import (
	"context"
	"ishari-backend/internal/core/entity"
)

type VerseUseCase interface {
	Create(ctx context.Context, input CreateVerseInput) (*entity.Verse, error)
	List(ctx context.Context, params ListParams) (*PaginatedResult[entity.Verse], error)
	Update(ctx context.Context, id uint, input UpdateVerseInput) (*entity.Verse, error)
	Delete(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	GetById(ctx context.Context, id uint) (*entity.Verse, error)
}

type CreateVerseInput struct {
	ChapterID       uint    `json:"chapter_id" gorm:"not null"`
	VerseNumber     uint    `json:"verse_number" gorm:"not null"`
	ArabicText      string  `json:"arabic_text" gorm:"type:text;not null"`
	Transliteration *string `json:"transliteration,omitempty" gorm:"type:text"`
}

type UpdateVerseInput struct {
	ChapterID       *uint   `json:"chapter_id" gorm:"not null"`
	VerseNumber     *uint   `json:"verse_number" gorm:"not null"`
	ArabicText      *string `json:"arabic_text" gorm:"type:text;not null"`
	Transliteration *string `json:"transliteration,omitempty" gorm:"type:text"`
}

type ListParams struct {
	Page            uint
	Limit           uint
	Search          string
	ChapterID       *uint
	ArabicText      string
	Transliteration string
}
