package dto

// CreateTranslationRequest struct for creating a new translation
type CreateTranslationRequest struct {
	VerseID         uint   `json:"verse_id" validate:"required"`
	LanguageCode    string `json:"language_code" validate:"required"`
	TranslationText string `json:"translation_text" validate:"required"`
	TranslatorName  string `json:"translator_name" validate:"required"`
}

// ListTranslationResponse struct for listing translations
type ListTranslationResponse struct {
	ID              uint    `json:"id"`
	VerseID         uint    `json:"verse_id"`
	LanguageCode    string  `json:"language_code"`
	TranslationText string  `json:"translation_text"`
	TranslatorName  *string `json:"translator_name"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// UpdateTranslationRequest struct for updating a translation
type UpdateTranslationRequest struct {
	VerseID         *uint   `json:"verse_id" validate:"required"`
	LanguageCode    *string `json:"language_code" validate:"required"`
	TranslationText *string `json:"translation_text" validate:"required"`
	TranslatorName  *string `json:"translator_name" validate:"required"`
}
