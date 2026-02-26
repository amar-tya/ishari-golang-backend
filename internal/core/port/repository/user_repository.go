package repository

import (
	"context"

	"ishari-backend/internal/core/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entity.User, error)
	UpdateLastLoginAt(ctx context.Context, userID uint) error
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, user *entity.User) error
	ListUsers(ctx context.Context, offset, limit int, search string) ([]entity.User, int64, error)
	BulkDelete(ctx context.Context, ids []uint) error
}
