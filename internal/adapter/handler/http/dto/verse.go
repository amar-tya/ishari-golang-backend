package dto

// CreateVerseRequest struct for creating a new verse
type CreateVerseRequest struct {
	ChapterID       uint    `json:"chapter_id"`
	VerseNumber     uint    `json:"verse_number"`
	ArabicText      string  `json:"arabic_text"`
	Transliteration *string `json:"transliteration,omitempty"`
}

// ListVerseResponse struct for listing verses
type ListVerseResponse struct {
	ID              uint    `json:"id"`
	ChapterID       uint    `json:"chapter_id"`
	VerseNumber     uint    `json:"verse_number"`
	ArabicText      string  `json:"arabic_text"`
	Transliteration *string `json:"transliteration,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// UpdateVerseRequest struct for updating a verse
type UpdateVerseRequest struct {
	ChapterID       *uint   `json:"chapter_id"`
	VerseNumber     *uint   `json:"verse_number"`
	ArabicText      *string `json:"arabic_text"`
	Transliteration *string `json:"transliteration,omitempty"`
}
