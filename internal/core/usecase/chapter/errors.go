package chapter

import "ishari-backend/internal/core/domain"

// Chapter-specific domain errors
var (
	// ErrChapterNotFound indicates the chapter does not exist
	ErrChapterNotFound = domain.NewNotFoundError("chapter not found", nil)

	// ErrBookNotFound indicates the book does not exist
	ErrBookNotFound = domain.NewNotFoundError("book not found", nil)

	// ErrInvalidChapterNumber indicates chapter number is invalid
	ErrInvalidChapterNumber = domain.NewInvalidInputError("chapter number must be greater than 0", nil)

	// ErrInvalidTitle indicates title is invalid
	ErrInvalidTitle = domain.NewInvalidInputError("title is required", nil)

	// ErrInvalidCategory indicates category is invalid
	ErrInvalidCategory = domain.NewInvalidInputError("category is required", nil)

	// ErrInvalidTotalVerses indicates total verses is invalid
	ErrInvalidTotalVerses = domain.NewInvalidInputError("total verses must be greater than 0", nil)
)
