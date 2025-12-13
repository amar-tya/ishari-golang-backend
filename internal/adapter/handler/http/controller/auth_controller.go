package controller

import (
	"time"

	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/adapter/handler/http/middleware"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/pkg/logger"
	"ishari-backend/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication HTTP requests
type AuthController struct {
	authUseCase portuc.AuthUseCase
	validate    validation.Validator
	log         logger.Logger
}

// NewAuthController creates a new AuthController instance
func NewAuthController(authUC portuc.AuthUseCase, v validation.Validator, l logger.Logger) *AuthController {
	return &AuthController{
		authUseCase: authUC,
		validate:    v,
		log:         l,
	}
}

// Login handles user authentication
// POST /api/v1/auth/login
func (h *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("Login body parse error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("Login validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	input := portuc.LoginInput{
		UsernameOrEmail: req.UsernameOrEmail,
		Password:        req.Password,
	}

	result, err := h.authUseCase.Login(c.UserContext(), input)
	if err != nil {
		h.log.Error("Login failed", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Build user response
	userResp := dto.UserResponse{
		ID:        result.User.ID,
		Username:  result.User.Username,
		Email:     result.User.Email,
		IsActive:  result.User.IsActive,
		CreatedAt: result.User.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: result.User.UpdatedAt.UTC().Format(time.RFC3339),
	}

	if result.User.LastLoginAt != nil {
		formatted := result.User.LastLoginAt.UTC().Format(time.RFC3339)
		userResp.LastLoginAt = &formatted
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "login successful",
		"data": dto.AuthResponse{
			User:         userResp,
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresAt:    result.ExpiresAt.UTC().Format(time.RFC3339),
		},
	})
}

// Logout invalidates the user's token
// POST /api/v1/auth/logout
func (h *AuthController) Logout(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "unauthorized",
		})
	}

	token := middleware.GetTokenFromContext(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "token not found",
		})
	}

	if err := h.authUseCase.Logout(c.UserContext(), claims.UserID, token); err != nil {
		h.log.Error("Logout failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "logout successful",
	})
}

// RefreshToken generates new tokens from refresh token
// POST /api/v1/auth/refresh
func (h *AuthController) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("RefreshToken body parse error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("RefreshToken validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	result, err := h.authUseCase.RefreshToken(c.UserContext(), req.RefreshToken)
	if err != nil {
		h.log.Error("RefreshToken failed", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Build user response
	userResp := dto.UserResponse{
		ID:        result.User.ID,
		Username:  result.User.Username,
		Email:     result.User.Email,
		IsActive:  result.User.IsActive,
		CreatedAt: result.User.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: result.User.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "token refreshed",
		"data": dto.AuthResponse{
			User:         userResp,
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresAt:    result.ExpiresAt.UTC().Format(time.RFC3339),
		},
	})
}

// GetMe returns current authenticated user info
// GET /api/v1/auth/me
func (h *AuthController) GetMe(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"email":    claims.Email,
		},
	})
}
