package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/factory"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/internal/core/usecase/auth"
)

// MockUserRepository for testing
type MockUserRepository struct {
	GetByUsernameOrEmailFunc func(ctx context.Context, usernameOrEmail string) (*entity.User, error)
	GetByIDFunc              func(ctx context.Context, id uint) (*entity.User, error)
	UpdateLastLoginAtFunc    func(ctx context.Context, userID uint) error
}

func (m *MockUserRepository) Create(ctx context.Context, u *entity.User) error { return nil }
func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *MockUserRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
	if m.GetByUsernameOrEmailFunc != nil {
		return m.GetByUsernameOrEmailFunc(ctx, usernameOrEmail)
	}
	return nil, errors.New("not found")
}
func (m *MockUserRepository) UpdateLastLoginAt(ctx context.Context, userID uint) error {
	if m.UpdateLastLoginAtFunc != nil {
		return m.UpdateLastLoginAtFunc(ctx, userID)
	}
	return nil
}
func (m *MockUserRepository) Delete(ctx context.Context, id uint) error        { return nil }
func (m *MockUserRepository) Update(ctx context.Context, u *entity.User) error { return nil }
func (m *MockUserRepository) ListUsers(ctx context.Context, o, l int, s string) ([]entity.User, int64, error) {
	return nil, 0, nil
}

// MockTokenService for testing
type MockTokenService struct {
	GenerateAccessTokenFunc  func(userID uint, username, email string) (string, time.Time, error)
	GenerateRefreshTokenFunc func(userID uint) (string, time.Time, error)
	ValidateAccessTokenFunc  func(token string) (*auth.Claims, error)
	ValidateRefreshTokenFunc func(token string) (*auth.RefreshClaims, error)
}

func (m *MockTokenService) GenerateAccessToken(userID uint, username, email string) (string, time.Time, error) {
	if m.GenerateAccessTokenFunc != nil {
		return m.GenerateAccessTokenFunc(userID, username, email)
	}
	return "access_token", time.Now().Add(15 * time.Minute), nil
}
func (m *MockTokenService) GenerateRefreshToken(userID uint) (string, time.Time, error) {
	if m.GenerateRefreshTokenFunc != nil {
		return m.GenerateRefreshTokenFunc(userID)
	}
	return "refresh_token", time.Now().Add(7 * 24 * time.Hour), nil
}
func (m *MockTokenService) ValidateAccessToken(token string) (*auth.Claims, error) {
	if m.ValidateAccessTokenFunc != nil {
		return m.ValidateAccessTokenFunc(token)
	}
	return nil, auth.ErrInvalidToken
}
func (m *MockTokenService) ValidateRefreshToken(token string) (*auth.RefreshClaims, error) {
	if m.ValidateRefreshTokenFunc != nil {
		return m.ValidateRefreshTokenFunc(token)
	}
	return nil, auth.ErrRefreshTokenInvalid
}

// MockBlacklist for testing
type MockBlacklist struct {
	tokens map[string]bool
}

func NewMockBlacklist() *MockBlacklist {
	return &MockBlacklist{tokens: make(map[string]bool)}
}
func (m *MockBlacklist) Add(token string, expiry time.Time) error {
	m.tokens[token] = true
	return nil
}
func (m *MockBlacklist) IsBlacklisted(token string) bool {
	return m.tokens[token]
}

// MockPasswordHasher for testing
type MockPasswordHasher struct{}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	return "hashed_" + password, nil
}
func (m *MockPasswordHasher) Compare(hashedPassword, password string) error {
	if hashedPassword == "hashed_"+password {
		return nil
	}
	return errors.New("password mismatch")
}

// Tests

func TestAuthUseCase_Login_Success(t *testing.T) {
	user := factory.NewUserFactory().
		WithID(1).
		WithUsername("testuser").
		WithEmail("test@example.com").
		WithPasswordHash("hashed_password123").
		WithIsActive(true).
		Build()

	mockRepo := &MockUserRepository{
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			return user, nil
		},
		UpdateLastLoginAtFunc: func(ctx context.Context, userID uint) error { return nil },
	}
	mockTokenService := &MockTokenService{}
	mockBlacklist := NewMockBlacklist()
	mockHasher := &MockPasswordHasher{}

	uc := auth.NewAuthUseCase(mockRepo, mockTokenService, mockBlacklist, mockHasher)
	input := portuc.LoginInput{
		UsernameOrEmail: "testuser",
		Password:        "password123",
	}

	result, err := uc.Login(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.AccessToken == "" {
		t.Error("expected access token, got empty")
	}
	if result.RefreshToken == "" {
		t.Error("expected refresh token, got empty")
	}
	if result.User.ID != 1 {
		t.Errorf("expected user ID 1, got %d", result.User.ID)
	}
}

func TestAuthUseCase_Login_InvalidCredentials(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			return nil, errors.New("not found")
		},
	}
	mockTokenService := &MockTokenService{}
	mockBlacklist := NewMockBlacklist()
	mockHasher := &MockPasswordHasher{}

	uc := auth.NewAuthUseCase(mockRepo, mockTokenService, mockBlacklist, mockHasher)
	input := portuc.LoginInput{
		UsernameOrEmail: "nonexistent",
		Password:        "password123",
	}

	_, err := uc.Login(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUseCase_Login_InactiveUser(t *testing.T) {
	user := factory.NewUserFactory().WithIsActive(false).WithPasswordHash("hashed_password123").Build()

	mockRepo := &MockUserRepository{
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			return user, nil
		},
	}
	mockTokenService := &MockTokenService{}
	mockBlacklist := NewMockBlacklist()
	mockHasher := &MockPasswordHasher{}

	uc := auth.NewAuthUseCase(mockRepo, mockTokenService, mockBlacklist, mockHasher)
	input := portuc.LoginInput{
		UsernameOrEmail: "testuser",
		Password:        "password123",
	}

	_, err := uc.Login(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, auth.ErrUserInactive) {
		t.Errorf("expected ErrUserInactive, got %v", err)
	}
}

func TestAuthUseCase_Logout_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	mockTokenService := &MockTokenService{
		ValidateAccessTokenFunc: func(token string) (*auth.Claims, error) {
			return &auth.Claims{
				UserID:   1,
				Username: "testuser",
				Email:    "test@example.com",
				Exp:      time.Now().Add(15 * time.Minute),
			}, nil
		},
	}
	mockBlacklist := NewMockBlacklist()
	mockHasher := &MockPasswordHasher{}

	uc := auth.NewAuthUseCase(mockRepo, mockTokenService, mockBlacklist, mockHasher)

	err := uc.Logout(context.Background(), 1, "valid_token")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !mockBlacklist.IsBlacklisted("valid_token") {
		t.Error("expected token to be blacklisted")
	}
}

func TestAuthUseCase_ValidateToken_Blacklisted(t *testing.T) {
	mockRepo := &MockUserRepository{}
	mockTokenService := &MockTokenService{}
	mockBlacklist := NewMockBlacklist()
	mockBlacklist.Add("blacklisted_token", time.Now().Add(time.Hour))
	mockHasher := &MockPasswordHasher{}

	uc := auth.NewAuthUseCase(mockRepo, mockTokenService, mockBlacklist, mockHasher)

	_, err := uc.ValidateToken(context.Background(), "blacklisted_token")

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, auth.ErrTokenBlacklisted) {
		t.Errorf("expected ErrTokenBlacklisted, got %v", err)
	}
}

func TestAuthUseCase_RefreshToken_Success(t *testing.T) {
	user := factory.NewUserFactory().WithID(1).WithIsActive(true).Build()

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*entity.User, error) {
			return user, nil
		},
	}
	mockTokenService := &MockTokenService{
		ValidateRefreshTokenFunc: func(token string) (*auth.RefreshClaims, error) {
			return &auth.RefreshClaims{
				UserID: 1,
				Exp:    time.Now().Add(7 * 24 * time.Hour),
			}, nil
		},
	}
	mockBlacklist := NewMockBlacklist()
	mockHasher := &MockPasswordHasher{}

	uc := auth.NewAuthUseCase(mockRepo, mockTokenService, mockBlacklist, mockHasher)

	result, err := uc.RefreshToken(context.Background(), "valid_refresh_token")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.AccessToken == "" {
		t.Error("expected new access token")
	}
	if result.RefreshToken == "" {
		t.Error("expected new refresh token")
	}
}
