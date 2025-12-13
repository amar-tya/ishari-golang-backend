package errors

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"ishari-backend/pkg/logger"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

func BadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message, "")
}

func NotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, "")
}

func InternalServerError(message string, detail string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, detail)
}

func Unauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, "")
}

func Forbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, "")
}

func Handler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var appErr *AppError
		if errors.As(err, &appErr) {
			return c.Status(appErr.Code).JSON(appErr)
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"code":    fiberErr.Code,
				"message": fiberErr.Message,
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
			"detail":  err.Error(),
		})
	}
}

// HandlerWithLogger is like Handler but also logs the error using provided logger.
func HandlerWithLogger(l logger.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Log the error with basic request context
		if l != nil && err != nil {
			l.Error("http error",
				"method", c.Method(),
				"path", c.Path(),
				"error", err.Error(),
			)
		}

		var appErr *AppError
		if errors.As(err, &appErr) {
			return c.Status(appErr.Code).JSON(appErr)
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"code":    fiberErr.Code,
				"message": fiberErr.Message,
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
			"detail":  err.Error(),
		})
	}
}
