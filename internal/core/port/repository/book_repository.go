package repository

import (
	"context"

	"ishari-backend/internal/core/entity"
)

type BookRepository interface {
	List(ctx context.Context, offset, limit int, search string) ([]entity.Book, int64, error)
	Create(ctx context.Context, book *entity.Book) error
	Edit(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*entity.Book, error)
}
