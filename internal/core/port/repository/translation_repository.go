package repository

import (
	"context"
	"ishari-backend/internal/core/entity"
)

type TranslationListFilter struct {
	Offset         uint
	Limit          uint
	Search         string
	VerseID        uint
	TranslatorName string
	LanguageCode   string
}

type TranslationRepository interface {
	Create(ctx context.Context, translation *entity.Translation) error
	List(ctx context.Context, filter TranslationListFilter) ([]entity.Translation, uint, error)
	Update(ctx context.Context, translation *entity.Translation) error
	Delete(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (*entity.Translation, error)
	GetByVerseId(ctx context.Context, verseId uint) ([]entity.Translation, error)
	GetDropdownData(ctx context.Context) (verses []entity.Verse, translatorNames []string, languageCodes []string, err error)
	BulkDelete(ctx context.Context, ids []uint) error
}
