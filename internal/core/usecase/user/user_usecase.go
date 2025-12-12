package user

import (
	"context"
	"regexp"
	"strings"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
)

// userUseCase implements the UserUseCase interface with business logic
type userUseCase struct {
	userRepo repository.UserRepository
	hasher   PasswordHasher
}

// NewUserUseCase creates a new UserUseCase instance
func NewUserUseCase(repo repository.UserRepository, hasher PasswordHasher) portuc.UserUseCase {
	return &userUseCase{
		userRepo: repo,
		hasher:   hasher,
	}
}

// Register creates a new user account
func (uc *userUseCase) Register(ctx context.Context, input portuc.RegisterUserInput) (*entity.User, error) {
	// Validate input
	if err := uc.validateRegistration(input); err != nil {
		return nil, err
	}

	// Check if email already exists
	existingByEmail, _ := uc.userRepo.GetByUsernameOrEmail(ctx, input.Email)
	if existingByEmail != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Check if username already exists
	existingByUsername, _ := uc.userRepo.GetByUsernameOrEmail(ctx, input.Username)
	if existingByUsername != nil {
		return nil, ErrUsernameAlreadyExists
	}

	// Hash password
	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user := &entity.User{
		Username:     input.Username,
		Email:        strings.ToLower(input.Email),
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	// Persist user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates user and returns user data
func (uc *userUseCase) Login(ctx context.Context, input portuc.LoginInput) (*entity.User, error) {
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

	// Update last login time
	_ = uc.userRepo.UpdateLastLoginAt(ctx, user.ID)

	return user, nil
}

// GetByID retrieves a user by their ID
func (uc *userUseCase) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Update modifies user information
func (uc *userUseCase) Update(ctx context.Context, id uint, input portuc.UpdateUserInput) (*entity.User, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Update fields if provided
	if input.Username != nil {
		if err := uc.validateUsername(*input.Username); err != nil {
			return nil, err
		}
		// Check if new username is already taken by another user
		existing, _ := uc.userRepo.GetByUsernameOrEmail(ctx, *input.Username)
		if existing != nil && existing.ID != id {
			return nil, ErrUsernameAlreadyExists
		}
		user.Username = *input.Username
	}

	if input.Email != nil {
		if err := uc.validateEmail(*input.Email); err != nil {
			return nil, err
		}
		// Check if new email is already taken by another user
		existing, _ := uc.userRepo.GetByUsernameOrEmail(ctx, *input.Email)
		if existing != nil && existing.ID != id {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = strings.ToLower(*input.Email)
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	// Persist changes
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes a user (soft delete)
func (uc *userUseCase) Delete(ctx context.Context, id uint) error {
	// Verify user exists
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	return uc.userRepo.Delete(ctx, id)
}

// List returns paginated users with optional search
func (uc *userUseCase) List(ctx context.Context, params portuc.ListUserParams) (*portuc.PaginatedResult[entity.User], error) {
	// Apply defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	offset := (params.Page - 1) * params.Limit

	users, total, err := uc.userRepo.ListUsers(ctx, offset, params.Limit, params.Search)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}

	return &portuc.PaginatedResult[entity.User]{
		Data:       users,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// validateRegistration validates registration input
func (uc *userUseCase) validateRegistration(input portuc.RegisterUserInput) error {
	if err := uc.validateUsername(input.Username); err != nil {
		return err
	}
	if err := uc.validateEmail(input.Email); err != nil {
		return err
	}
	if err := uc.validatePassword(input.Password); err != nil {
		return err
	}
	return nil
}

// validateUsername checks username format
func (uc *userUseCase) validateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return ErrInvalidUsername
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if !matched {
		return ErrInvalidUsername
	}
	return nil
}

// validateEmail checks email format
func (uc *userUseCase) validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// validatePassword checks password requirements
func (uc *userUseCase) validatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}
	return nil
}
