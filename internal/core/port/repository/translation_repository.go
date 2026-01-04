package repository

import (
	"context"
	"ishari-backend/internal/core/entity"
)

type TranslationRepository interface {
	Create(ctx context.Context, translation *entity.Translation) error
	List(ctx context.Context, offset, limit uint, search string) ([]entity.Translation, uint, error)
	Update(ctx context.Context, translation *entity.Translation) error
	Delete(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (*entity.Translation, error)
	GetByVerseId(ctx context.Context, verseId uint) ([]entity.Translation, error)
}
