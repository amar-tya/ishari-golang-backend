package dto

// CreateTranslationRequest struct for creating a new translation
type CreateTranslationRequest struct {
	VerseID         uint   `json:"verse_id" validate:"required"`
	LanguageCode    string `json:"language_code" validate:"required"`
	TranslationText string `json:"translation_text" validate:"required"`
	TranslatorName  string `json:"translator_name" validate:"required"`
}

// BulkDeleteTranslationRequest represents the HTTP request for bulk deleting translations
type BulkDeleteTranslationRequest struct {
	IDs []uint `json:"ids" query:"ids" form:"ids" validate:"required,min=1"`
}

// VerseDropdownItem is a minimal verse representation for dropdowns
type VerseDropdownItem struct {
	ID         uint   `json:"id"`
	ArabicText string `json:"arabic_text"`
}

// ListTranslationResponse struct for listing translations
type ListTranslationResponse struct {
	ID              uint               `json:"id"`
	VerseID         uint               `json:"verse_id"`
	Verse           *VerseDropdownItem `json:"verse,omitempty"`
	LanguageCode    string             `json:"language_code"`
	TranslationText string             `json:"translation_text"`
	TranslatorName  *string            `json:"translator_name"`
	CreatedAt       string             `json:"created_at"`
	UpdatedAt       string             `json:"updated_at"`
}

// TranslationDropdownResponse is the response for GET /translations/dropdown
type TranslationDropdownResponse struct {
	Verses          []VerseDropdownItem `json:"verses"`
	TranslatorNames []string            `json:"translator_names"`
	LanguageCodes   []string            `json:"language_codes"`
}

// UpdateTranslationRequest struct for updating a translation
type UpdateTranslationRequest struct {
	VerseID         *uint   `json:"verse_id" validate:"required"`
	LanguageCode    *string `json:"language_code" validate:"required"`
	TranslationText *string `json:"translation_text" validate:"required"`
	TranslatorName  *string `json:"translator_name" validate:"required"`
}
