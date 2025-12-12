package repository

import (
	"context"

	"ishari-backend/internal/core/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entity.User, error)
	UpdateLastLoginAt(ctx context.Context, userID uint) error
}
