package repository

import (
	"context"
	"ishari-backend/internal/core/entity"
)

type VerseRepository interface {
	Create(ctx context.Context, verse *entity.Verse) error
	List(ctx context.Context, offset, limit uint, search string) ([]entity.Verse, uint, error)
	Update(ctx context.Context, verse *entity.Verse) error
	Delete(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	GetById(ctx context.Context, id uint) (*entity.Verse, error)
}
