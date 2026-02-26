package repository

import (
	"context"

	"ishari-backend/internal/core/entity"
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, bookmark *entity.Bookmark) error
	GetBookmarkByID(ctx context.Context, id uint) (*entity.Bookmark, error)
	GetBookmarkByUserIDAndVerseID(ctx context.Context, userID uint, verseID uint) (*entity.Bookmark, error)
	ListBookmarksByUserID(ctx context.Context, userID uint, offset, limit int, sort string) ([]entity.Bookmark, int64, error)
	UpdateBookmark(ctx context.Context, bookmark *entity.Bookmark) error
	DeleteBookmark(ctx context.Context, id uint) error
}
