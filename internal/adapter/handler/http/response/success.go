package response

import "github.com/gofiber/fiber/v2"

// SuccessResponse represents a standard success response structure
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse represents a paginated response structure
type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	Count      int `json:"count"`
}

// SendSuccess sends a standardized success response
func SendSuccess(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	response := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	return ctx.Status(statusCode).JSON(response)
}

// SendOK sends a 200 OK success response
func SendOK(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(SuccessResponse{
		Status: "success",
		Data:   data,
	})
}

// SendCreated sends a 201 Created success response
func SendCreated(ctx *fiber.Ctx, message string, data interface{}) error {
	return SendSuccess(ctx, fiber.StatusCreated, message, data)
}

// SendPaginated sends a paginated response
func SendPaginated(ctx *fiber.Ctx, data interface{}, page, limit int, total int64, totalPages, count int) error {
	response := PaginatedResponse{
		Data: data,
		Meta: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
			Count:      count,
		},
	}
	return ctx.JSON(response)
}
