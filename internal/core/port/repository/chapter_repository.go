package repository

import (
	"context"

	"ishari-backend/internal/core/entity"
)

type ChapterRepository interface {
	CreateChapter(ctx context.Context, chapter *entity.Chapter) error
	ListChapters(ctx context.Context, offset, limit int, search string) ([]entity.Chapter, int64, error)
	GetChaptersByBookID(ctx context.Context, bookID uint) ([]entity.Chapter, int64, error)
	GetChapterByID(ctx context.Context, id uint) (*entity.Chapter, error)
	UpdateChapter(ctx context.Context, chapter *entity.Chapter) error
	DeleteChapter(ctx context.Context, id uint) error
	DeleteChapters(ctx context.Context, ids []uint) error
}
