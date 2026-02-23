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

// NewUserRepository creates a new UserRepository implementation
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by their primary key
func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsernameOrEmail retrieves a user by username or email
func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLoginAt updates the last login timestamp for a user
func (r *userRepository) UpdateLastLoginAt(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("last_login_at", time.Now()).Error
}

// Delete performs a soft delete on the user
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

// Update modifies an existing user in the database
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).
		Model(user).
		Updates(map[string]interface{}{
			"username":  user.Username,
			"email":     user.Email,
			"is_active": user.IsActive,
		}).Error
}

// ListUsers retrieves paginated users with optional search
func (r *userRepository) ListUsers(ctx context.Context, offset, limit int, search string) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.User{}).Where("deleted_at IS NULL")

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// BulkDelete performs a soft delete on multiple users
func (r *userRepository) BulkDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id IN ?", ids).
		Update("deleted_at", time.Now()).Error
}
