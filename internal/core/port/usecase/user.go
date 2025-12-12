package usecase

import (
	"context"

	"ishari-backend/internal/core/entity"
)

// UserUseCase defines the interface for user business logic operations.
// This interface is defined in the port layer following Dependency Inversion Principle.
type UserUseCase interface {
	// Register creates a new user account
	Register(ctx context.Context, input RegisterUserInput) (*entity.User, error)

	// Login authenticates user and returns user data
	Login(ctx context.Context, input LoginInput) (*entity.User, error)

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uint) (*entity.User, error)

	// Update modifies user information
	Update(ctx context.Context, id uint, input UpdateUserInput) (*entity.User, error)

	// Delete removes a user (soft delete)
	Delete(ctx context.Context, id uint) error

	// List returns paginated users with optional search
	List(ctx context.Context, params ListUserParams) (*PaginatedResult[entity.User], error)
}

// RegisterUserInput contains data required for user registration
type RegisterUserInput struct {
	Username string
	Email    string
	Password string
}

// LoginInput contains data required for user login
type LoginInput struct {
	UsernameOrEmail string
	Password        string
}

// UpdateUserInput contains data for updating user profile
type UpdateUserInput struct {
	Username *string
	Email    *string
	IsActive *bool
}

// ListUserParams contains pagination and filter parameters
type ListUserParams struct {
	Page   int
	Limit  int
	Search string
}

// PaginatedResult is a generic pagination wrapper
type PaginatedResult[T any] struct {
	Data       []T
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}
