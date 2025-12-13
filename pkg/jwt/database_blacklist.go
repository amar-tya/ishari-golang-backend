package jwt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
)

// DatabaseBlacklist implements TokenBlacklist using database storage
type DatabaseBlacklist struct {
	repo repository.RefreshTokenRepository
}

// NewDatabaseBlacklist creates a new database-backed blacklist
func NewDatabaseBlacklist(repo repository.RefreshTokenRepository) *DatabaseBlacklist {
	bl := &DatabaseBlacklist{repo: repo}

	// Start periodic cleanup of expired tokens
	go bl.periodicCleanup()

	return bl
}

// hashToken creates a SHA256 hash of the token
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Add adds a token to the blacklist (marks it as revoked)
func (bl *DatabaseBlacklist) Add(token string, expiry time.Time) error {
	tokenHash := hashToken(token)
	return bl.repo.RevokeByTokenHash(context.Background(), tokenHash)
}

// IsBlacklisted checks if a token has been revoked
func (bl *DatabaseBlacklist) IsBlacklisted(token string) bool {
	tokenHash := hashToken(token)
	refreshToken, err := bl.repo.GetByTokenHash(context.Background(), tokenHash)
	if err != nil {
		// Token not found in DB - could be an access token (not stored)
		return false
	}
	return refreshToken.IsRevoked()
}

// StoreRefreshToken stores a new refresh token in the database
func (bl *DatabaseBlacklist) StoreRefreshToken(userID uint, token string, expiresAt time.Time, userAgent, ipAddress string) error {
	tokenHash := hashToken(token)
	refreshToken := &entity.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}
	return bl.repo.Create(context.Background(), refreshToken)
}

// ValidateRefreshToken checks if a refresh token is valid (exists and not revoked)
func (bl *DatabaseBlacklist) ValidateRefreshToken(token string) (*entity.RefreshToken, error) {
	tokenHash := hashToken(token)
	refreshToken, err := bl.repo.GetByTokenHash(context.Background(), tokenHash)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user (logout all devices)
func (bl *DatabaseBlacklist) RevokeAllUserTokens(userID uint) error {
	return bl.repo.RevokeAllByUserID(context.Background(), userID)
}

// periodicCleanup removes expired tokens periodically
func (bl *DatabaseBlacklist) periodicCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		_ = bl.repo.DeleteExpired(context.Background())
	}
}
