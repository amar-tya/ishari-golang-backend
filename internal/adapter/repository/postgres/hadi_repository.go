package postgres

import (
	"context"

	"ishari-backend/internal/core/entity"
	portrepo "ishari-backend/internal/core/port/repository"

	"gorm.io/gorm"
)

type hadiRepository struct {
	db *gorm.DB
}

// NewHadiRepository creates a new hadi repository instance
func NewHadiRepository(db *gorm.DB) portrepo.HadiRepository {
	return &hadiRepository{db: db}
}

func (r *hadiRepository) Create(ctx context.Context, hadi *entity.Hadi) error {
	return r.db.WithContext(ctx).Create(hadi).Error
}

func (r *hadiRepository) GetByID(ctx context.Context, id int) (*entity.Hadi, error) {
	var hadi entity.Hadi
	err := r.db.WithContext(ctx).First(&hadi, id).Error
	if err != nil {
		return nil, err
	}
	return &hadi, nil
}

func (r *hadiRepository) List(ctx context.Context, limit, offset int) ([]entity.Hadi, int64, error) {
	var hadis []entity.Hadi
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Hadi{})

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Order("created_at desc").Find(&hadis).Error
	if err != nil {
		return nil, 0, err
	}

	return hadis, total, nil
}

func (r *hadiRepository) Update(ctx context.Context, hadi *entity.Hadi) error {
	return r.db.WithContext(ctx).Save(hadi).Error
}

func (r *hadiRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&entity.Hadi{}, id).Error
}
