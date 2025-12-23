package response

import (
	"errors"
	"ishari-backend/internal/core/domain"
	"ishari-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents a standard error response structure
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SendError sends a standardized error response
func SendError(ctx *fiber.Ctx, statusCode int, message string, err error, log logger.Logger, logContext string) error {
	if log != nil && logContext != "" {
		log.Error(logContext, "error", err)
	}

	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return ctx.Status(statusCode).JSON(response)
}

// SendBadRequest sends a 400 Bad Request error
func SendBadRequest(ctx *fiber.Ctx, message string, err error, log logger.Logger, logContext string) error {
	return SendError(ctx, fiber.StatusBadRequest, message, err, log, logContext)
}

// SendUnauthorized sends a 401 Unauthorized error
func SendUnauthorized(ctx *fiber.Ctx, message string, err error, log logger.Logger, logContext string) error {
	return SendError(ctx, fiber.StatusUnauthorized, message, err, log, logContext)
}

// SendNotFound sends a 404 Not Found error
func SendNotFound(ctx *fiber.Ctx, message string, err error, log logger.Logger, logContext string) error {
	return SendError(ctx, fiber.StatusNotFound, message, err, log, logContext)
}

// SendInternalError sends a 500 Internal Server Error
func SendInternalError(ctx *fiber.Ctx, message string, err error, log logger.Logger, logContext string) error {
	return SendError(ctx, fiber.StatusInternalServerError, message, err, log, logContext)
}

// SendValidationError sends a 400 Bad Request for validation errors
func SendValidationError(ctx *fiber.Ctx, err error, log logger.Logger, logContext string) error {
	return SendBadRequest(ctx, "validation failed", err, log, logContext)
}

// SendParseError sends a 400 Bad Request for body parsing errors
func SendParseError(ctx *fiber.Ctx, err error, log logger.Logger, logContext string) error {
	return SendBadRequest(ctx, "invalid request body", err, log, logContext)
}

// SendDomainError maps domain errors to appropriate HTTP responses
// This function automatically determines the HTTP status code based on the domain error type
func SendDomainError(ctx *fiber.Ctx, err error, log logger.Logger) error {
	var domainErr *domain.DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Type {
		case domain.ErrTypeNotFound:
			return SendNotFound(ctx, domainErr.Message, err, log, "")
		case domain.ErrTypeInvalidInput:
			return SendBadRequest(ctx, domainErr.Message, err, log, "")
		case domain.ErrTypeUnauthorized:
			return SendUnauthorized(ctx, domainErr.Message, err, log, "")
		case domain.ErrTypeConflict:
			return SendError(ctx, fiber.StatusConflict, domainErr.Message, err, log, "")
		case domain.ErrTypeInternal:
			return SendInternalError(ctx, domainErr.Message, err, log, "")
		default:
			return SendInternalError(ctx, domainErr.Message, err, log, "")
		}
	}
	// Fallback for non-domain errors
	return SendInternalError(ctx, "internal server error", err, log, "")
}
