package auth

import (
	"context"

	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
	userusecase "ishari-backend/internal/core/usecase/user"
)

// authUseCase implements AuthUseCase interface
type authUseCase struct {
	userRepo     repository.UserRepository
	tokenService TokenService
	blacklist    TokenBlacklist
	hasher       userusecase.PasswordHasher
}

// NewAuthUseCase creates a new AuthUseCase instance
func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenService TokenService,
	blacklist TokenBlacklist,
	hasher userusecase.PasswordHasher,
) portuc.AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
		blacklist:    blacklist,
		hasher:       hasher,
	}
}

// Login authenticates user and returns token pair
func (uc *authUseCase) Login(ctx context.Context, input portuc.LoginInput) (*portuc.AuthResult, error) {
	// Find user by username or email
	user, err := uc.userRepo.GetByUsernameOrEmail(ctx, input.UsernameOrEmail)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := uc.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, expiresAt, err := uc.tokenService.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Update last login time
	_ = uc.userRepo.UpdateLastLoginAt(ctx, user.ID)

	return &portuc.AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout invalidates the user's token
func (uc *authUseCase) Logout(ctx context.Context, userID uint, token string) error {
	// Validate token first
	claims, err := uc.tokenService.ValidateAccessToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	// Check if token belongs to the user
	if claims.UserID != userID {
		return ErrInvalidToken
	}

	// Add to blacklist
	return uc.blacklist.Add(token, claims.Exp)
}

// RefreshToken generates new access token from refresh token
func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*portuc.AuthResult, error) {
	// Validate refresh token
	claims, err := uc.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// Get user from database
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Generate new tokens
	accessToken, expiresAt, err := uc.tokenService.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	newRefreshToken, _, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &portuc.AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// ValidateToken validates token and returns claims
func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*portuc.TokenClaims, error) {
	// Check blacklist first
	if uc.blacklist.IsBlacklisted(token) {
		return nil, ErrTokenBlacklisted
	}

	// Validate token
	claims, err := uc.tokenService.ValidateAccessToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &portuc.TokenClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
	}, nil
}

// GetUserFromContext is a helper to get user info from validated token
func GetUserFromContext(ctx context.Context) (*portuc.TokenClaims, bool) {
	claims, ok := ctx.Value("user").(*portuc.TokenClaims)
	return claims, ok
}
