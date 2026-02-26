package translation

import "ishari-backend/internal/core/domain"

var (
	ErrTranslationNotFound = domain.NewNotFoundError("translation not found", nil)

	ErrInvalidTranslationText = domain.NewInvalidInputError("translation text is required", nil)
	ErrInvalidLanguageCode = domain.NewInvalidInputError("language code is required", nil)
)
