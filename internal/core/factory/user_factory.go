package factory

import (
	"fmt"
	"time"

	"ishari-backend/internal/core/entity"
)

// UserFactory provides a fluent API for creating User entities for testing
type UserFactory struct {
	id           uint
	username     string
	email        string
	passwordHash string
	isActive     bool
	lastLoginAt  *time.Time
	createdAt    time.Time
	updatedAt    time.Time
	counter      int
}

// NewUserFactory creates a new UserFactory with sensible defaults
func NewUserFactory() *UserFactory {
	now := time.Now()
	return &UserFactory{
		id:           0,
		username:     "",
		email:        "",
		passwordHash: "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4L9kJEZvP8..vKiy", // "password123"
		isActive:     true,
		lastLoginAt:  nil,
		createdAt:    now,
		updatedAt:    now,
		counter:      0,
	}
}

// WithID sets a specific ID
func (f *UserFactory) WithID(id uint) *UserFactory {
	f.id = id
	return f
}

// WithUsername sets a specific username
func (f *UserFactory) WithUsername(username string) *UserFactory {
	f.username = username
	return f
}

// WithEmail sets a specific email
func (f *UserFactory) WithEmail(email string) *UserFactory {
	f.email = email
	return f
}

// WithPasswordHash sets a specific password hash
func (f *UserFactory) WithPasswordHash(hash string) *UserFactory {
	f.passwordHash = hash
	return f
}

// WithIsActive sets the active status
func (f *UserFactory) WithIsActive(active bool) *UserFactory {
	f.isActive = active
	return f
}

// WithLastLoginAt sets the last login time
func (f *UserFactory) WithLastLoginAt(t time.Time) *UserFactory {
	f.lastLoginAt = &t
	return f
}

// WithCreatedAt sets the created time
func (f *UserFactory) WithCreatedAt(t time.Time) *UserFactory {
	f.createdAt = t
	return f
}

// WithUpdatedAt sets the updated time
func (f *UserFactory) WithUpdatedAt(t time.Time) *UserFactory {
	f.updatedAt = t
	return f
}

// Build creates a new User entity with the configured values
func (f *UserFactory) Build() *entity.User {
	f.counter++

	username := f.username
	if username == "" {
		username = fmt.Sprintf("testuser%d", f.counter)
	}

	email := f.email
	if email == "" {
		email = fmt.Sprintf("test%d@example.com", f.counter)
	}

	id := f.id
	if id == 0 {
		id = uint(f.counter)
	}

	return &entity.User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: f.passwordHash,
		IsActive:     f.isActive,
		LastLoginAt:  f.lastLoginAt,
		CreatedAt:    f.createdAt,
		UpdatedAt:    f.updatedAt,
	}
}

// BuildList creates multiple User entities
func (f *UserFactory) BuildList(count int) []entity.User {
	users := make([]entity.User, count)
	for i := 0; i < count; i++ {
		users[i] = *f.Build()
	}
	return users
}

// Reset resets the factory to default values
func (f *UserFactory) Reset() *UserFactory {
	return NewUserFactory()
}
