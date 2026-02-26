package dto

// CreateHadiRequest defines the payload for creating a new hadi
type CreateHadiRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty"`
	ImageURL    *string `json:"image_url" validate:"omitempty,url"`
}

// UpdateHadiRequest defines the payload for updating an existing hadi
type UpdateHadiRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty"`
	ImageURL    *string `json:"image_url" validate:"omitempty,url"`
}
