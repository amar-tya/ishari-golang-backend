package postgres

import (
	"context"
	"strings"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"

	"gorm.io/gorm"
)

type chapterRepository struct {
	db *gorm.DB
}

// NewChapterRepository creates a new chapter repository instance
func NewChapterRepository(db *gorm.DB) repository.ChapterRepository {
	return &chapterRepository{db: db}
}

// CreateChapter creates a new chapter in the database
func (r *chapterRepository) CreateChapter(ctx context.Context, chapter *entity.Chapter) error {
	return r.db.WithContext(ctx).Create(chapter).Error
}

// ListChapters retrieves paginated chapters with optional search
func (r *chapterRepository) ListChapters(ctx context.Context, offset, limit int, search string) ([]entity.Chapter, int64, error) {
	var (
		total    int64
		chapters []entity.Chapter
	)

	base := r.db.WithContext(ctx).Model(&entity.Chapter{})
	if search = strings.TrimSpace(search); search != "" {
		q := "%" + search + "%"
		base = base.Where("title ILIKE ? OR category ILIKE ?", q, q)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := base.Order("book_id ASC, chapter_number ASC").Offset(offset).Limit(limit)
	if err := query.Find(&chapters).Error; err != nil {
		return nil, 0, err
	}

	return chapters, total, nil
}

// GetChaptersByBookID retrieves all chapters for a specific book
func (r *chapterRepository) GetChaptersByBookID(ctx context.Context, bookID uint) ([]entity.Chapter, int64, error) {
	var (
		total    int64
		chapters []entity.Chapter
	)

	base := r.db.WithContext(ctx).Model(&entity.Chapter{}).Where("book_id = ?", bookID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := base.Order("chapter_number ASC").Find(&chapters).Error; err != nil {
		return nil, 0, err
	}

	return chapters, total, nil
}

// GetChapterByID retrieves a chapter by its ID
func (r *chapterRepository) GetChapterByID(ctx context.Context, id uint) (*entity.Chapter, error) {
	var chapter entity.Chapter
	if err := r.db.WithContext(ctx).First(&chapter, id).Error; err != nil {
		return nil, err
	}
	return &chapter, nil
}

// UpdateChapter updates an existing chapter
func (r *chapterRepository) UpdateChapter(ctx context.Context, chapter *entity.Chapter) error {
	return r.db.WithContext(ctx).Save(chapter).Error
}

// DeleteChapter removes a chapter by ID
func (r *chapterRepository) DeleteChapter(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Chapter{}, id).Error
}
