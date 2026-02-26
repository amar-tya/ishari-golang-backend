package repository

import (
	"context"

	"ishari-backend/internal/core/entity"
)

type HadiRepository interface {
	Create(ctx context.Context, hadi *entity.Hadi) error
	GetByID(ctx context.Context, id int) (*entity.Hadi, error)
	List(ctx context.Context, limit, offset int) ([]entity.Hadi, int64, error)
	Update(ctx context.Context, hadi *entity.Hadi) error
	Delete(ctx context.Context, id int) error
}
