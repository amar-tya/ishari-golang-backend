package domain

import "fmt"

// ErrorType represents the category of domain error
type ErrorType string

const (
	// ErrTypeNotFound indicates a resource was not found
	ErrTypeNotFound ErrorType = "NOT_FOUND"

	// ErrTypeInvalidInput indicates invalid input data
	ErrTypeInvalidInput ErrorType = "INVALID_INPUT"

	// ErrTypeUnauthorized indicates unauthorized access
	ErrTypeUnauthorized ErrorType = "UNAUTHORIZED"

	// ErrTypeConflict indicates a conflict with existing data
	ErrTypeConflict ErrorType = "CONFLICT"

	// ErrTypeInternal indicates an internal system error
	ErrTypeInternal ErrorType = "INTERNAL"
)

// DomainError represents a domain-level error with classification
type DomainError struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewNotFoundError creates a new NOT_FOUND error
func NewNotFoundError(message string, err error) *DomainError {
	return &DomainError{
		Type:    ErrTypeNotFound,
		Message: message,
		Err:     err,
	}
}

// NewInvalidInputError creates a new INVALID_INPUT error
func NewInvalidInputError(message string, err error) *DomainError {
	return &DomainError{
		Type:    ErrTypeInvalidInput,
		Message: message,
		Err:     err,
	}
}

// NewUnauthorizedError creates a new UNAUTHORIZED error
func NewUnauthorizedError(message string, err error) *DomainError {
	return &DomainError{
		Type:    ErrTypeUnauthorized,
		Message: message,
		Err:     err,
	}
}

// NewConflictError creates a new CONFLICT error
func NewConflictError(message string, err error) *DomainError {
	return &DomainError{
		Type:    ErrTypeConflict,
		Message: message,
		Err:     err,
	}
}

// NewInternalError creates a new INTERNAL error
func NewInternalError(message string, err error) *DomainError {
	return &DomainError{
		Type:    ErrTypeInternal,
		Message: message,
		Err:     err,
	}
}
