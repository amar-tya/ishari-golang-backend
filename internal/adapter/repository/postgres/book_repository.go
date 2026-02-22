package postgres

import (
	"context"
	"strings"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"

	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) repository.BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) List(ctx context.Context, offset, limit int, search string) ([]entity.Book, int64, error) {
	var (
		total int64
		books []entity.Book
	)

	base := r.db.WithContext(ctx).Model(&entity.Book{}).Where("deleted_at IS NULL")
	if search = strings.TrimSpace(search); search != "" {
		q := "%" + search + "%"
		base = base.Where("title ILIKE ? OR author ILIKE ?", q, q)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := base.Order("created_at DESC").Offset(offset).Limit(limit)
	if err := query.Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *bookRepository) Create(ctx context.Context, book *entity.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepository) Edit(ctx context.Context, book *entity.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *bookRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&entity.Book{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *bookRepository) GetById(ctx context.Context, id int64) (*entity.Book, error) {
	var book entity.Book
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}
