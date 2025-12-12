package controller

import (
	"math"
	"strconv"
	"time"

	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/core/usecase"
	"ishari-backend/pkg/logger"
	"ishari-backend/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

type BookController struct {
	bookUseCase usecase.BookUseCase
	validate    validation.Validator
	log         logger.Logger
}

func NewBookController(bookUseCase usecase.BookUseCase, v validation.Validator, l logger.Logger) *BookController {
	return &BookController{bookUseCase: bookUseCase, validate: v, log: l}
}

func (h *BookController) CreateBook(c *fiber.Ctx) error {
	var req dto.CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("CreateBook body parse error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("CreateBook validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	in := usecase.CreateBookInput{
		Title:         req.Title,
		Author:        req.Author,
		Description:   req.Description,
		PublishedYear: req.PublishedYear,
		CoverImageURL: req.CoverImageURL,
	}

	book, err := h.bookUseCase.CreateBook(c.UserContext(), in)
	if err != nil {
		h.log.Error("CreateBook failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	resp := dto.BookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		Description:   book.Description,
		PublishedYear: book.PublishedYear,
		CoverImageURL: book.CoverImageURL,
		CreatedAt:     book.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     book.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *BookController) ListBooks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")

	books, total, err := h.bookUseCase.ListBooks(c.UserContext(), page, limit, search)
	if err != nil {
		h.log.Error("ListBooks failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to list books",
			"error":   err.Error(),
		})
	}

	var totalPages int
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	out := make([]dto.BookResponse, 0, len(books))
	for _, b := range books {
		out = append(out, dto.BookResponse{
			ID:            b.ID,
			Title:         b.Title,
			Author:        b.Author,
			Description:   b.Description,
			PublishedYear: b.PublishedYear,
			CoverImageURL: b.CoverImageURL,
			CreatedAt:     b.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:     b.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return c.JSON(fiber.Map{
		"data": out,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"count":       len(books),
		},
	})
}

func (h *BookController) EditBook(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var req dto.CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("EditBook body parse error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("EditBook validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	in := usecase.CreateBookInput{
		Title:         req.Title,
		Author:        req.Author,
		Description:   req.Description,
		PublishedYear: req.PublishedYear,
		CoverImageURL: req.CoverImageURL,
	}

	book, err := h.bookUseCase.EditBook(c.UserContext(), int64(id), in)
	if err != nil {
		h.log.Error("EditBook failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	resp := dto.BookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		Description:   book.Description,
		PublishedYear: book.PublishedYear,
		CoverImageURL: book.CoverImageURL,
		CreatedAt:     book.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     book.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *BookController) DeleteBook(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.bookUseCase.DeleteBook(c.UserContext(), int64(id)); err != nil {
		h.log.Error("DeleteBook failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "book deleted successfully",
	})
}

func (h *BookController) GetBookById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	book, err := h.bookUseCase.GetBookById(c.UserContext(), int64(id))
	if err != nil {
		h.log.Error("GetBookById failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	resp := dto.BookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		Description:   book.Description,
		PublishedYear: book.PublishedYear,
		CoverImageURL: book.CoverImageURL,
		CreatedAt:     book.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     book.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
