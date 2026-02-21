package usecase

import (
	"context"
	"ishari-backend/internal/core/entity"
)

type TranslationUseCase interface {
	Create(ctx context.Context, input CreateTranslationInput) (*entity.Translation, error)
	List(ctx context.Context, params ListParams) (*PaginatedResult[entity.Translation], error)
	Update(ctx context.Context, id uint, input UpdateTranslationInput) (*entity.Translation, error)
	Delete(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (*entity.Translation, error)
	GetByVerseId(ctx context.Context, verseId uint) ([]entity.Translation, error)
	GetDropdownData(ctx context.Context) (*TranslationDropdownData, error)
	BulkDelete(ctx context.Context, ids []uint) error
}

type TranslationDropdownData struct {
	Verses          []entity.Verse `json:"verses"`
	TranslatorNames []string       `json:"translator_names"`
	LanguageCodes   []string       `json:"language_codes"`
}

type CreateTranslationInput struct {
	VerseID         uint    `json:"verse_id" gorm:"not null"`
	LanguageCode    string  `json:"language_code" gorm:"type:varchar(10);not null"`
	TranslationText string  `json:"translation_text" gorm:"type:text;not null"`
	TranslatorName  *string `json:"translator_name,omitempty" gorm:"type:varchar(255)"`
}

type UpdateTranslationInput struct {
	VerseID         *uint   `json:"verse_id" gorm:"not null"`
	LanguageCode    *string `json:"language_code" gorm:"type:varchar(10);not null"`
	TranslationText *string `json:"translation_text" gorm:"type:text;not null"`
	TranslatorName  *string `json:"translator_name,omitempty" gorm:"type:varchar(255)"`
}
