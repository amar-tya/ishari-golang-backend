package postgres

import (
	"context"
	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
	"strings"

	"gorm.io/gorm"
)

type TranslationRepository struct {
	db *gorm.DB
}

// NewTranslationRepository creates a new translation repository implementation
func NewTranslationRepository(db *gorm.DB) repository.TranslationRepository {
	return &TranslationRepository{db: db}
}

// Create implements TranslationRepository.
func (r *TranslationRepository) Create(ctx context.Context, translation *entity.Translation) error {
	return r.db.WithContext(ctx).Create(translation).Error
}

// List implements TranslationRepository.
func (r *TranslationRepository) List(ctx context.Context, offset, limit uint, search string) ([]entity.Translation, uint, error) {
	var (
		total        int64
		translations []entity.Translation
	)

	base := r.db.WithContext(ctx).Model(&entity.Translation{})
	if search = strings.TrimSpace(search); search != "" {
		q := "%" + search + "%"
		base = base.Where("translation_text ILIKE ? OR translator_name ILIKE ?", q, q)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := base.Order("verse_id ASC, language_code ASC").Offset(int(offset)).Limit(int(limit))
	if err := query.Find(&translations).Error; err != nil {
		return nil, 0, err
	}

	return translations, uint(total), nil
}

// GetByVerseId implements TranslationRepository.
func (r *TranslationRepository) GetByVerseId(ctx context.Context, verseId uint) ([]entity.Translation, error) {
	var translations []entity.Translation
	if err := r.db.WithContext(ctx).Where("verse_id = ?", verseId).Find(&translations).Error; err != nil {
		return nil, err
	}
	return translations, nil
}

// GetById implements TranslationRepository.
func (r *TranslationRepository) GetById(ctx context.Context, id uint) (*entity.Translation, error) {
	var translation entity.Translation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&translation).Error; err != nil {
		return nil, err
	}
	return &translation, nil
}

// Update implements TranslationRepository.
func (r *TranslationRepository) Update(ctx context.Context, translation *entity.Translation) error {
	return r.db.WithContext(ctx).Save(translation).Error
}

// Delete implements TranslationRepository.
func (r *TranslationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Translation{}, id).Error
}
