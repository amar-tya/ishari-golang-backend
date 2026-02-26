package jwt

import (
	"errors"
	"time"

	"ishari-backend/internal/core/usecase/auth"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService implements TokenService using golang-jwt
type JWTService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTService creates a new JWTService
func NewJWTService(secret string, accessTTL, refreshTTL time.Duration) *JWTService {
	return &JWTService{
		secretKey:       []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// accessTokenClaims for access tokens
type accessTokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// refreshTokenClaims for refresh tokens
type refreshTokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new access token
func (s *JWTService) GenerateAccessToken(userID uint, username, email, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.accessTokenTTL)

	claims := accessTokenClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// GenerateRefreshToken creates a new refresh token
func (s *JWTService) GenerateRefreshToken(userID uint) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.refreshTokenTTL)

	claims := refreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// ValidateAccessToken validates and parses access token
func (s *JWTService) ValidateAccessToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, auth.ErrTokenExpired
		}
		return nil, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok || !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	// Verify it's an access token
	if claims.Subject != "access" {
		return nil, auth.ErrInvalidToken
	}

	return &auth.Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
		Exp:      claims.ExpiresAt.Time,
	}, nil
}

// ValidateRefreshToken validates and parses refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*auth.RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, auth.ErrTokenExpired
		}
		return nil, auth.ErrRefreshTokenInvalid
	}

	claims, ok := token.Claims.(*refreshTokenClaims)
	if !ok || !token.Valid {
		return nil, auth.ErrRefreshTokenInvalid
	}

	// Verify it's a refresh token
	if claims.Subject != "refresh" {
		return nil, auth.ErrRefreshTokenInvalid
	}

	return &auth.RefreshClaims{
		UserID: claims.UserID,
		Exp:    claims.ExpiresAt.Time,
	}, nil
}
