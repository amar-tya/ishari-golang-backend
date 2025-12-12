package postgres

import (
	"context"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateLastLoginAt(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("last_login_at", time.Now()).Error
}
