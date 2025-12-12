package dto

type CreateBookRequest struct {
	Title         string  `json:"title" validate:"required,min=1,max=255"`
	Author        *string `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	Description   *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	PublishedYear *int    `json:"published_year,omitempty" validate:"omitempty,min=0"`
	CoverImageURL *string `json:"cover_image_url,omitempty" validate:"omitempty,uri"`
}

type BookResponse struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Author        *string `json:"author,omitempty"`
	Description   *string `json:"description,omitempty"`
	PublishedYear *int    `json:"published_year,omitempty"`
	CoverImageURL *string `json:"cover_image_url,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
