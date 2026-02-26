package usecase

import (
	"context"
	"time"

	"ishari-backend/internal/core/entity"
)

// AuthUseCase defines the interface for authentication operations
type AuthUseCase interface {
	// Login authenticates user and returns token pair
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)

	// Logout invalidates the user's token
	Logout(ctx context.Context, userID uint, token string) error

	// RefreshToken generates new access token from refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error)

	// ValidateToken validates token and returns claims
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
}

// AuthResult contains tokens and user info after successful auth
type AuthResult struct {
	User         *entity.User
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// TokenClaims contains decoded JWT claims
type TokenClaims struct {
	UserID   uint
	Username string
	Email    string
	Role     string
}
