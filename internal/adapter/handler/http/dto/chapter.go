package dto

// ListChapterResponse represent the HTTP response for listing chapters
type ListChapterResponse struct {
	ID            uint          `json:"id"`
	BookID        uint          `json:"book_id"`
	ChapterNumber uint          `json:"chapter_number"`
	Title         string        `json:"title"`
	Category      string        `json:"category"`
	Description   *string       `json:"description,omitempty"`
	TotalVerses   uint          `json:"total_verses"`
	Book          *BookResponse `json:"book,omitempty"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
}

// CreateChapterRequest represents the HTTP request for creating a chapter
type CreateChapterRequest struct {
	BookID        uint    `json:"book_id" validate:"required"`
	ChapterNumber uint    `json:"chapter_number" validate:"required"`
	Title         string  `json:"title" validate:"required"`
	Category      string  `json:"category" validate:"required"`
	Description   *string `json:"description"`
	TotalVerses   uint    `json:"total_verses" validate:"required"`
}

// UpdateChapterRequest represents the HTTP request for updating a chapter
type UpdateChapterRequest struct {
	BookID        *uint   `json:"book_id"`
	ChapterNumber *uint   `json:"chapter_number"`
	Title         *string `json:"title"`
	Category      *string `json:"category"`
	Description   *string `json:"description"`
	TotalVerses   *uint   `json:"total_verses"`
}

// BulkDeleteChapterRequest represents the HTTP request for bulk deleting chapters
type BulkDeleteChapterRequest struct {
	IDs []uint `json:"ids" query:"ids" form:"ids" validate:"required,min=1"`
}
