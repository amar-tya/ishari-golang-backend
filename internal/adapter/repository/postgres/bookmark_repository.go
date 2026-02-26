package postgres

import (
	"context"
	"errors"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"

	"gorm.io/gorm"
)

type bookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository creates a new instance of BookmarkRepository
func NewBookmarkRepository(db *gorm.DB) repository.BookmarkRepository {
	return &bookmarkRepository{db: db}
}

// CreateBookmark creates a new bookmark
func (r *bookmarkRepository) CreateBookmark(ctx context.Context, bookmark *entity.Bookmark) error {
	return r.db.WithContext(ctx).Create(bookmark).Error
}

// GetBookmarkByID retrieves a bookmark by its ID
func (r *bookmarkRepository) GetBookmarkByID(ctx context.Context, id uint) (*entity.Bookmark, error) {
	var bookmark entity.Bookmark
	err := r.db.WithContext(ctx).First(&bookmark, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &bookmark, nil
}

// GetBookmarkByUserIDAndVerseID retrieves a bookmark by UserID and VerseID
func (r *bookmarkRepository) GetBookmarkByUserIDAndVerseID(ctx context.Context, userID uint, verseID uint) (*entity.Bookmark, error) {
	var bookmark entity.Bookmark
	err := r.db.WithContext(ctx).Where("user_id = ? AND verse_id = ?", userID, verseID).First(&bookmark).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, to allow checking for existence
		}
		return nil, err
	}
	return &bookmark, nil
}

// ListBookmarksByUserID retrieves all bookmarks for a specific user with pagination
func (r *bookmarkRepository) ListBookmarksByUserID(ctx context.Context, userID uint, offset, limit int, sort string) ([]entity.Bookmark, int64, error) {
	var bookmarks []entity.Bookmark
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Bookmark{}).Where("user_id = ?", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}

	err := query.Offset(offset).Limit(limit).Find(&bookmarks).Error
	if err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

// UpdateBookmark updates an existing bookmark
func (r *bookmarkRepository) UpdateBookmark(ctx context.Context, bookmark *entity.Bookmark) error {
	return r.db.WithContext(ctx).Save(bookmark).Error
}

// DeleteBookmark deletes a bookmark by its ID
func (r *bookmarkRepository) DeleteBookmark(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Bookmark{}, id).Error
}
