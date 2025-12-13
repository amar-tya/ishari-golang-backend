package usecase

import (
	"context"
	"ishari-backend/internal/core/entity"
)

type BookUseCase interface {
	ListBooks(ctx context.Context, page, limit int, search string) ([]entity.Book, int64, error)
	CreateBook(ctx context.Context, input CreateBookInput) (*entity.Book, error)
	EditBook(ctx context.Context, id int64, input CreateBookInput) (*entity.Book, error)
	DeleteBook(ctx context.Context, id int64) error
	GetBookById(ctx context.Context, id int64) (*entity.Book, error)
}

type CreateBookInput struct {
	Title         string
	Author        *string
	Description   *string
	PublishedYear *int
	CoverImageURL *string
}
