package repository

import (
	"context"
	"ishari-backend/internal/core/entity"
)

type VerseFilter struct {
	Offset          uint
	Limit           uint
	Search          string
	ChapterID       *uint
	ArabicText      string
	Transliteration string
}

type VerseRepository interface {
	Create(ctx context.Context, verse *entity.Verse) error
	List(ctx context.Context, filter VerseFilter) ([]entity.Verse, uint, error)
	Update(ctx context.Context, verse *entity.Verse) error
	Delete(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	GetById(ctx context.Context, id uint) (*entity.Verse, error)
}
