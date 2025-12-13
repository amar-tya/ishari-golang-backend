package auth

import "errors"

// Domain-specific errors for authentication operations
var (
	// ErrInvalidToken indicates the token is malformed or invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired indicates the token has expired
	ErrTokenExpired = errors.New("token has expired")

	// ErrTokenBlacklisted indicates the token has been invalidated/logged out
	ErrTokenBlacklisted = errors.New("token has been invalidated")

	// ErrInvalidCredentials indicates wrong username/email or password
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserInactive indicates the user account is deactivated
	ErrUserInactive = errors.New("user account is inactive")

	// ErrRefreshTokenInvalid indicates refresh token is invalid
	ErrRefreshTokenInvalid = errors.New("invalid refresh token")
)
