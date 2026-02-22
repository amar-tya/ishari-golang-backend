package user

import "errors"

// Domain-specific errors for user operations.
// These errors are part of the business logic and are independent of infrastructure.
var (
	// ErrEmailAlreadyExists indicates the email is already registered
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrUsernameAlreadyExists indicates the username is already taken
	ErrUsernameAlreadyExists = errors.New("username already exists")

	// ErrInvalidCredentials indicates wrong username/email or password
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound indicates the user does not exist
	ErrUserNotFound = errors.New("user not found")

	// ErrUserInactive indicates the user account is deactivated
	ErrUserInactive = errors.New("user account is inactive")

	// ErrInvalidPassword indicates password does not meet requirements
	ErrInvalidPassword = errors.New("password must be at least 8 characters")

	// ErrInvalidEmail indicates email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrInvalidUsername indicates username does not meet requirements
	ErrInvalidUsername = errors.New("username must be 3-50 characters, alphanumeric and underscores only")

	// ErrUnauthorizedOperation indicates the user is not allowed to perform the action
	ErrUnauthorizedOperation = errors.New("unauthorized operation")
)
