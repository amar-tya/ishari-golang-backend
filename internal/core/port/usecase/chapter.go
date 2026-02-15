package usecase

import (
	"context"

	"ishari-backend/internal/core/entity"
)

type ChapterUsecase interface {
	Create(ctx context.Context, input CreateChapterInput) (*entity.Chapter, error)
	List(ctx context.Context, params ListChapterInput) (*PaginatedResult[entity.Chapter], error)
	GetByBookID(ctx context.Context, bookID uint) (*PaginatedResult[entity.Chapter], error)
	GetByID(ctx context.Context, id uint) (*entity.Chapter, error)
	Update(ctx context.Context, id uint, input UpdateChapterInput) (*entity.Chapter, error)
	Delete(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
}

// CreateChapterInput contains data required to create a new chapter.
type CreateChapterInput struct {
	BookID        uint
	ChapterNumber uint
	Title         string
	Category      string
	Description   *string
	TotalVerses   uint
}

// UpdateChapterInput contains data required to update a chapter.
type UpdateChapterInput struct {
	BookID        *uint
	ChapterNumber *uint
	Title         *string
	Category      *string
	Description   *string
	TotalVerses   *uint
}

// ListChapterInput contains data required to list chapters.
type ListChapterInput struct {
	Page   int
	Limit  int
	Search string
}

// PaginatedResult is a generic pagination wrapper
// use user.go for reference
