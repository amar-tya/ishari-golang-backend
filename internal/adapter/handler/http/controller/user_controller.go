package controller

import (
	"math"
	"strconv"
	"time"

	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/core/entity"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/pkg/logger"
	"ishari-backend/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	userUseCase portuc.UserUseCase
	validate    validation.Validator
	log         logger.Logger
}

// NewUserController creates a new UserController instance
func NewUserController(userUseCase portuc.UserUseCase, v validation.Validator, l logger.Logger) *UserController {
	return &UserController{
		userUseCase: userUseCase,
		validate:    v,
		log:         l,
	}
}

// Register handles user registration
// POST /api/v1/users/register
func (h *UserController) Register(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("Register body parse error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("Register validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	input := portuc.RegisterUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.userUseCase.Register(c.UserContext(), input)
	if err != nil {
		h.log.Error("Register failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "user registered successfully",
		"data":    h.toUserResponse(user),
	})
}

// Login handles user authentication
// POST /api/v1/users/login
func (h *UserController) Login(c *fiber.Ctx) error {
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

	user, err := h.userUseCase.Login(c.UserContext(), input)
	if err != nil {
		h.log.Error("Login failed", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "login successful",
		"data":    h.toUserResponse(user),
	})
}

// GetUserByID retrieves a user by ID
// GET /api/v1/users/:id
func (h *UserController) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid user id",
		})
	}

	user, err := h.userUseCase.GetByID(c.UserContext(), uint(id))
	if err != nil {
		h.log.Error("GetUserByID failed", "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   h.toUserResponse(user),
	})
}

// UpdateUser modifies user information
// PUT /api/v1/users/:id
func (h *UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid user id",
		})
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("UpdateUser body parse error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("UpdateUser validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	input := portuc.UpdateUserInput{
		Username: req.Username,
		Email:    req.Email,
		IsActive: req.IsActive,
	}

	user, err := h.userUseCase.Update(c.UserContext(), uint(id), input)
	if err != nil {
		h.log.Error("UpdateUser failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "user updated successfully",
		"data":    h.toUserResponse(user),
	})
}

// DeleteUser removes a user
// DELETE /api/v1/users/:id
func (h *UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid user id",
		})
	}

	if err := h.userUseCase.Delete(c.UserContext(), uint(id)); err != nil {
		h.log.Error("DeleteUser failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "user deleted successfully",
	})
}

// ListUsers retrieves paginated users
// GET /api/v1/users
func (h *UserController) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")

	params := portuc.ListUserParams{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	result, err := h.userUseCase.List(c.UserContext(), params)
	if err != nil {
		h.log.Error("ListUsers failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to list users",
			"error":   err.Error(),
		})
	}

	var totalPages int
	if limit > 0 {
		totalPages = int(math.Ceil(float64(result.Total) / float64(limit)))
	}

	out := make([]dto.UserResponse, 0, len(result.Data))
	for _, u := range result.Data {
		out = append(out, h.toUserResponse(&u))
	}

	return c.JSON(fiber.Map{
		"data": out,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       result.Total,
			"total_pages": totalPages,
			"count":       len(result.Data),
		},
	})
}

// toUserResponse converts entity.User to dto.UserResponse
func (h *UserController) toUserResponse(user *entity.User) dto.UserResponse {
	resp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
	}

	if user.LastLoginAt != nil {
		formatted := user.LastLoginAt.UTC().Format(time.RFC3339)
		resp.LastLoginAt = &formatted
	}

	return resp
}
