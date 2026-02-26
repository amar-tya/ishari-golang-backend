package auth

import (
	"time"
)

// TokenService defines the interface for JWT token operations
// This abstraction allows the usecase to be independent of the JWT implementation
type TokenService interface {
	// GenerateAccessToken creates a new access token
	GenerateAccessToken(userID uint, username, email, role string) (string, time.Time, error)

	// GenerateRefreshToken creates a new refresh token
	GenerateRefreshToken(userID uint) (string, time.Time, error)

	// ValidateAccessToken validates and parses access token
	ValidateAccessToken(token string) (*Claims, error)

	// ValidateRefreshToken validates and parses refresh token
	ValidateRefreshToken(token string) (*RefreshClaims, error)
}

// Claims represents the access token claims
type Claims struct {
	UserID   uint
	Username string
	Email    string
	Role     string
	Exp      time.Time
}

// RefreshClaims represents the refresh token claims
type RefreshClaims struct {
	UserID uint
	Exp    time.Time
}

// TokenBlacklist interface for tracking invalidated tokens
type TokenBlacklist interface {
	// Add adds a token to the blacklist
	Add(token string, expiry time.Time) error

	// IsBlacklisted checks if token is blacklisted
	IsBlacklisted(token string) bool
}
