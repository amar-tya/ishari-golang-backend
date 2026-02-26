package controller

import (
	"math"
	"strconv"

	"ishari-backend/internal/adapter/handler/http/dto"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/pkg/logger"
	"ishari-backend/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

type HadiController struct {
	hadiUseCase portuc.HadiUseCase
	validate    validation.Validator
	log         logger.Logger
}

func NewHadiController(hadiUseCase portuc.HadiUseCase, v validation.Validator, l logger.Logger) *HadiController {
	return &HadiController{
		hadiUseCase: hadiUseCase,
		validate:    v,
		log:         l,
	}
}

func (h *HadiController) Create(c *fiber.Ctx) error {
	var req dto.CreateHadiRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("CreateHadi body parse error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("CreateHadi validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	resp, err := h.hadiUseCase.Create(c.UserContext(), req)
	if err != nil {
		h.log.Error("CreateHadi failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": resp,
	})
}

func (h *HadiController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid id format",
		})
	}

	resp, err := h.hadiUseCase.GetByID(c.UserContext(), id)
	if err != nil {
		h.log.Error("GetHadiByID failed", "error", err)
		status := fiber.StatusInternalServerError
		if err.Error() == "hadi not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": resp,
	})
}

func (h *HadiController) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	resp, total, err := h.hadiUseCase.List(c.UserContext(), page, limit)
	if err != nil {
		h.log.Error("ListHadis failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	var totalPages int
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": resp,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"count":       len(resp),
		},
	})
}

func (h *HadiController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid id format",
		})
	}

	var req dto.UpdateHadiRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("UpdateHadi body parse error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("UpdateHadi validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	resp, err := h.hadiUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		h.log.Error("UpdateHadi failed", "error", err)
		status := fiber.StatusInternalServerError
		if err.Error() == "hadi not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": resp,
	})
}

func (h *HadiController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid id format",
		})
	}

	if err := h.hadiUseCase.Delete(c.UserContext(), id); err != nil {
		h.log.Error("DeleteHadi failed", "error", err)
		status := fiber.StatusInternalServerError
		if err.Error() == "hadi not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "hadi deleted successfully",
	})
}
