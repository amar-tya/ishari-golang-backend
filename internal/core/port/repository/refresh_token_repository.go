package repository

import (
	"context"

	"ishari-backend/internal/core/entity"
)

// RefreshTokenRepository defines the interface for refresh token persistence
type RefreshTokenRepository interface {
	// Create stores a new refresh token
	Create(ctx context.Context, token *entity.RefreshToken) error

	// GetByTokenHash retrieves a refresh token by its hash
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)

	// RevokeByTokenHash revokes a specific token
	RevokeByTokenHash(ctx context.Context, tokenHash string) error

	// RevokeAllByUserID revokes all tokens for a user (logout from all devices)
	RevokeAllByUserID(ctx context.Context, userID uint) error

	// DeleteExpired removes expired tokens from the database
	DeleteExpired(ctx context.Context) error
}
