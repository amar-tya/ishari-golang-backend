package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/factory"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/internal/core/usecase/user"
)

// MockUserRepository is a manual mock for testing
type MockUserRepository struct {
	CreateFunc               func(ctx context.Context, user *entity.User) error
	GetByIDFunc              func(ctx context.Context, id uint) (*entity.User, error)
	GetByUsernameOrEmailFunc func(ctx context.Context, usernameOrEmail string) (*entity.User, error)
	UpdateLastLoginAtFunc    func(ctx context.Context, userID uint) error
	DeleteFunc               func(ctx context.Context, id uint) error
	UpdateFunc               func(ctx context.Context, user *entity.User) error
	ListUsersFunc            func(ctx context.Context, offset, limit int, search string) ([]entity.User, int64, error)
}

func (m *MockUserRepository) Create(ctx context.Context, u *entity.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, u)
	}
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
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

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, u *entity.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, u)
	}
	return nil
}

func (m *MockUserRepository) ListUsers(ctx context.Context, offset, limit int, search string) ([]entity.User, int64, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(ctx, offset, limit, search)
	}
	return nil, 0, nil
}

// MockPasswordHasher is a manual mock for password hashing
type MockPasswordHasher struct {
	HashFunc    func(password string) (string, error)
	CompareFunc func(hashedPassword, password string) error
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	if m.HashFunc != nil {
		return m.HashFunc(password)
	}
	return "hashed_" + password, nil
}

func (m *MockPasswordHasher) Compare(hashedPassword, password string) error {
	if m.CompareFunc != nil {
		return m.CompareFunc(hashedPassword, password)
	}
	if hashedPassword == "hashed_"+password {
		return nil
	}
	return errors.New("password mismatch")
}

// Tests

func TestUserUseCase_Register_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			return nil, errors.New("not found")
		},
		CreateFunc: func(ctx context.Context, u *entity.User) error {
			u.ID = 1
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
			return nil
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	input := portuc.RegisterUserInput{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	result, err := uc.Register(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", result.Username)
	}
	if result.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", result.Email)
	}
}

func TestUserUseCase_Register_DuplicateEmail(t *testing.T) {
	existingUser := factory.NewUserFactory().WithEmail("test@example.com").Build()

	mockRepo := &MockUserRepository{
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			if usernameOrEmail == "test@example.com" {
				return existingUser, nil
			}
			return nil, errors.New("not found")
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	input := portuc.RegisterUserInput{
		Username: "newuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err := uc.Register(context.Background(), input)

	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
	if !errors.Is(err, user.ErrEmailAlreadyExists) {
		t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestUserUseCase_Register_InvalidPassword(t *testing.T) {
	mockRepo := &MockUserRepository{}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	input := portuc.RegisterUserInput{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "short", // less than 8 characters
	}

	_, err := uc.Register(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid password, got nil")
	}
	if !errors.Is(err, user.ErrInvalidPassword) {
		t.Errorf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestUserUseCase_Login_Success(t *testing.T) {
	existingUser := factory.NewUserFactory().
		WithID(1).
		WithUsername("testuser").
		WithEmail("test@example.com").
		WithPasswordHash("hashed_password123").
		WithIsActive(true).
		Build()

	mockRepo := &MockUserRepository{
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			return existingUser, nil
		},
		UpdateLastLoginAtFunc: func(ctx context.Context, userID uint) error {
			return nil
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	input := portuc.LoginInput{
		UsernameOrEmail: "testuser",
		Password:        "password123",
	}

	result, err := uc.Login(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestUserUseCase_Login_InvalidCredentials(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			return nil, errors.New("not found")
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	input := portuc.LoginInput{
		UsernameOrEmail: "nonexistent",
		Password:        "password123",
	}

	_, err := uc.Login(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid credentials, got nil")
	}
	if !errors.Is(err, user.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserUseCase_Login_InactiveUser(t *testing.T) {
	inactiveUser := factory.NewUserFactory().
		WithIsActive(false).
		WithPasswordHash("hashed_password123").
		Build()

	mockRepo := &MockUserRepository{
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			return inactiveUser, nil
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	input := portuc.LoginInput{
		UsernameOrEmail: "testuser",
		Password:        "password123",
	}

	_, err := uc.Login(context.Background(), input)

	if err == nil {
		t.Error("expected error for inactive user, got nil")
	}
	if !errors.Is(err, user.ErrUserInactive) {
		t.Errorf("expected ErrUserInactive, got %v", err)
	}
}

func TestUserUseCase_GetByID_Success(t *testing.T) {
	existingUser := factory.NewUserFactory().WithID(1).Build()

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*entity.User, error) {
			return existingUser, nil
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)

	result, err := uc.GetByID(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestUserUseCase_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*entity.User, error) {
			return nil, errors.New("not found")
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)

	_, err := uc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserUseCase_Update_Success(t *testing.T) {
	existingUser := factory.NewUserFactory().
		WithID(1).
		WithUsername("olduser").
		WithEmail("old@example.com").
		Build()

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*entity.User, error) {
			return existingUser, nil
		},
		GetByUsernameOrEmailFunc: func(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
			return nil, errors.New("not found")
		},
		UpdateFunc: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	newUsername := "newuser"
	input := portuc.UpdateUserInput{
		Username: &newUsername,
	}

	result, err := uc.Update(context.Background(), 1, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Username != "newuser" {
		t.Errorf("expected username 'newuser', got %s", result.Username)
	}
}

func TestUserUseCase_Delete_Success(t *testing.T) {
	existingUser := factory.NewUserFactory().WithID(1).Build()

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*entity.User, error) {
			return existingUser, nil
		},
		DeleteFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)

	// Create context with authorized user
	ctx := portuc.NewContextWithUser(context.Background(), &portuc.TokenClaims{UserID: 1})

	err := uc.Delete(ctx, 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUserUseCase_Delete_Unauthorized(t *testing.T) {
	mockRepo := &MockUserRepository{}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)

	// Case 1: No user in context
	err := uc.Delete(context.Background(), 1)
	if !errors.Is(err, user.ErrUnauthorizedOperation) {
		t.Errorf("expected ErrUnauthorizedOperation when no user in context, got %v", err)
	}

	// Case 2: Different user in context (IDOR attempt)
	ctx := portuc.NewContextWithUser(context.Background(), &portuc.TokenClaims{UserID: 2})
	err = uc.Delete(ctx, 1)
	if !errors.Is(err, user.ErrUnauthorizedOperation) {
		t.Errorf("expected ErrUnauthorizedOperation when deleting another user, got %v", err)
	}
}

func TestUserUseCase_Delete_NotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uint) (*entity.User, error) {
			return nil, errors.New("not found")
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)

	// Create context with authorized user
	ctx := portuc.NewContextWithUser(context.Background(), &portuc.TokenClaims{UserID: 999})

	err := uc.Delete(ctx, 999)

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserUseCase_List_Success(t *testing.T) {
	users := factory.NewUserFactory().BuildList(3)

	mockRepo := &MockUserRepository{
		ListUsersFunc: func(ctx context.Context, offset, limit int, search string) ([]entity.User, int64, error) {
			return users, 3, nil
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	params := portuc.ListUserParams{
		Page:  1,
		Limit: 10,
	}

	result, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Total)
	}
	if len(result.Data) != 3 {
		t.Errorf("expected 3 users, got %d", len(result.Data))
	}
}

func TestUserUseCase_List_WithDefaults(t *testing.T) {
	mockRepo := &MockUserRepository{
		ListUsersFunc: func(ctx context.Context, offset, limit int, search string) ([]entity.User, int64, error) {
			// Verify defaults are applied
			if limit != 20 {
				t.Errorf("expected default limit 20, got %d", limit)
			}
			if offset != 0 {
				t.Errorf("expected offset 0, got %d", offset)
			}
			return []entity.User{}, 0, nil
		},
	}
	mockHasher := &MockPasswordHasher{}

	uc := user.NewUserUseCase(mockRepo, mockHasher)
	params := portuc.ListUserParams{
		Page:  0, // Should default to 1
		Limit: 0, // Should default to 20
	}

	_, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
