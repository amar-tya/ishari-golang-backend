package dto

import "time"

// BookmarkResponse represents the HTTP response for a bookmark
type BookmarkResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	VerseID   uint      `json:"verse_id"`
	Note      *string   `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// ListBookmarkResponse represents the HTTP response for listing bookmarks
type ListBookmarkResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	VerseID   uint      `json:"verse_id"`
	Note      *string   `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateBookmarkRequest represents the HTTP request for creating a bookmark
type CreateBookmarkRequest struct {
	VerseID uint    `json:"verse_id" validate:"required"`
	Note    *string `json:"note"`
}

// UpdateBookmarkRequest represents the HTTP request for updating a bookmark
type UpdateBookmarkRequest struct {
	Note *string `json:"note"`
}
