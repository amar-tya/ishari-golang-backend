package verse

import "ishari-backend/internal/core/domain"

var (
	ErrVerseNotFound      = domain.NewNotFoundError("verse not found", nil)
	ErrChapterNotFound    = domain.NewNotFoundError("chapter not found", nil)
	ErrInvalidVerseNumber = domain.NewInvalidInputError("verse number must be greater than 0", nil)
	ErrInvalidVerseText   = domain.NewInvalidInputError("verse text is required", nil)
)
