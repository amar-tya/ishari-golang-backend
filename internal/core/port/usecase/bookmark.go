package usecase

import (
	"context"

	"ishari-backend/internal/core/entity"
)

// BookmarkInput structures

type CreateBookmarkInput struct {
	UserID  uint
	VerseID uint
	Note    *string
}

type UpdateBookmarkInput struct {
	UserID uint
	Note   *string
}

type ListBookmarkInput struct {
	UserID uint
	Page   int
	Limit  int
	Sort   string // Support optional sorting, e.g. "created_at asc", "created_at desc"
}

// BookmarkUsecase interface declaration
type BookmarkUsecase interface {
	Create(ctx context.Context, input CreateBookmarkInput) (*entity.Bookmark, error)
	GetByID(ctx context.Context, id uint, userID uint) (*entity.Bookmark, error)
	ListByUserID(ctx context.Context, input ListBookmarkInput) (*PaginatedResult[entity.Bookmark], error)
	Update(ctx context.Context, id uint, input UpdateBookmarkInput) (*entity.Bookmark, error)
	Delete(ctx context.Context, id uint, userID uint) error
}
