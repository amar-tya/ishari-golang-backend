package response

import (
	"time"

	"ishari-backend/internal/core/entity"
)

type HadiResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func MapHadiResponse(h *entity.Hadi) *HadiResponse {
	if h == nil {
		return nil
	}

	return &HadiResponse{
		ID:          h.ID,
		Name:        h.Name,
		Description: h.Description,
		ImageURL:    h.ImageURL,
		CreatedAt:   h.CreatedAt,
		UpdatedAt:   h.UpdatedAt,
	}
}

func MapHadiListResponse(hadis []entity.Hadi) []HadiResponse {
	responses := make([]HadiResponse, 0, len(hadis))
	for _, h := range hadis {
		responses = append(responses, *MapHadiResponse(&h))
	}
	return responses
}
