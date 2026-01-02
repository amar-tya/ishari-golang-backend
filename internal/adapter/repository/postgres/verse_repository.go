package postgres

import (
	"context"
	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
	"strings"

	"gorm.io/gorm"
)

type VerseRepository struct {
	db *gorm.DB
}

// NewVerseRepository creates a new verse repository implementation
func NewVerseRepository(db *gorm.DB) repository.VerseRepository {
	return &VerseRepository{db: db}
}

// Create implements VerseRepository.
func (r *VerseRepository) Create(ctx context.Context, verse *entity.Verse) error {
	return r.db.WithContext(ctx).Create(verse).Error
}

// List implements VerseRepository.
func (r *VerseRepository) List(ctx context.Context, offset, limit uint, search string) ([]entity.Verse, uint, error) {
	var (
		total  int64
		verses []entity.Verse
	)

	base := r.db.WithContext(ctx).Model(&entity.Verse{})
	if search = strings.TrimSpace(search); search != "" {
		q := "%" + search + "%"
		base = base.Where("arabic_text ILIKE ? OR transliteration ILIKE ?", q, q)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := base.Order("chapter_id ASC, verse_number ASC").Offset(int(offset)).Limit(int(limit))
	if err := query.Find(&verses).Error; err != nil {
		return nil, 0, err
	}

	return verses, uint(total), nil
}

// GetById implements VerseRepository.
func (r *VerseRepository) GetById(ctx context.Context, id uint) (*entity.Verse, error) {
	var verse entity.Verse
	if err := r.db.WithContext(ctx).First(&verse, id).Error; err != nil {
		return nil, err
	}
	return &verse, nil
}

// Update implements VerseRepository.
func (r *VerseRepository) Update(ctx context.Context, verse *entity.Verse) error {
	return r.db.WithContext(ctx).Save(verse).Error
}

// Delete implements VerseRepository.
func (r *VerseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Verse{}, id).Error
}
