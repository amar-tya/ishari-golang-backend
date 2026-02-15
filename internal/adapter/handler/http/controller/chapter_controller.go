package controller

import (
	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/adapter/handler/http/response"
	"ishari-backend/internal/core/entity"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/pkg/logger"
	"ishari-backend/pkg/validation"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ChapterController handles chapter-related HTTP requests
type ChapterController struct {
	chapterUsecase portuc.ChapterUsecase
	validate       validation.Validator
	log            logger.Logger
}

// NewChapterController creates a new chapter controller
func NewChapterController(chapterUsecase portuc.ChapterUsecase, v validation.Validator, l logger.Logger) *ChapterController {
	return &ChapterController{
		chapterUsecase: chapterUsecase,
		validate:       v,
		log:            l,
	}
}

// List handles chapter listing
// GET /api/v1/chapters
func (c *ChapterController) List(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))
	search := ctx.Query("search", "")

	params := portuc.ListChapterInput{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	result, err := c.chapterUsecase.List(ctx.UserContext(), params)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	var totalPages int
	if limit > 0 {
		totalPages = int(math.Ceil(float64(result.Total) / float64(limit)))
	}

	out := make([]dto.ListChapterResponse, 0, len(result.Data))
	for _, u := range result.Data {
		out = append(out, c.toListChapterResponse(&u))
	}

	return response.SendPaginated(ctx, out, page, limit, result.Total, totalPages, len(result.Data))
}

// toListChapterResponse converts a Chapter to a ListChapterResponse
func (c *ChapterController) toListChapterResponse(chapter *entity.Chapter) dto.ListChapterResponse {
	return dto.ListChapterResponse{
		ID:            chapter.ID,
		BookID:        chapter.BookID,
		ChapterNumber: chapter.ChapterNumber,
		Title:         chapter.Title,
		Category:      chapter.Category,
		Description:   chapter.Description,
		TotalVerses:   chapter.TotalVerses,
		CreatedAt:     chapter.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     chapter.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

// GetByID handles getting a chapter by ID
// GET /api/v1/chapters/:id
func (c *ChapterController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid chapter ID", err, nil, "")
	}

	chapter, err := c.chapterUsecase.GetByID(ctx.UserContext(), uint(id))
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, c.toListChapterResponse(chapter))
}

// Create handles creating a new chapter
// POST /api/v1/chapters
func (c *ChapterController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateChapterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.SendParseError(ctx, err, c.log, "Create chapter body parse error")
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Create chapter validation failed")
	}

	input := portuc.CreateChapterInput{
		BookID:        req.BookID,
		ChapterNumber: req.ChapterNumber,
		Title:         req.Title,
		Category:      req.Category,
		Description:   req.Description,
		TotalVerses:   req.TotalVerses,
	}

	chapter, err := c.chapterUsecase.Create(ctx.UserContext(), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendCreated(ctx, "chapter created successfully", c.toListChapterResponse(chapter))
}

// Update handles updating a chapter
// PUT /api/v1/chapters/:id
func (c *ChapterController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid chapter ID", err, nil, "")
	}

	var req dto.UpdateChapterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.SendParseError(ctx, err, c.log, "Update chapter body parse error")
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Update chapter validation failed")
	}

	input := portuc.UpdateChapterInput{
		BookID:        req.BookID,
		ChapterNumber: req.ChapterNumber,
		Title:         req.Title,
		Category:      req.Category,
		Description:   req.Description,
		TotalVerses:   req.TotalVerses,
	}

	chapter, err := c.chapterUsecase.Update(ctx.UserContext(), uint(id), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, c.toListChapterResponse(chapter))
}

// Delete handles deleting a chapter
// DELETE /api/v1/chapters/:id
func (c *ChapterController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid chapter ID", err, nil, "")
	}

	if err := c.chapterUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, fiber.Map{"message": "chapter deleted successfully"})
}

// BulkDelete handles bulk deleting chapters
// POST /api/v1/chapters/bulk-delete
func (c *ChapterController) BulkDelete(ctx *fiber.Ctx) error {
	var req dto.BulkDeleteChapterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.SendParseError(ctx, err, c.log, "Bulk delete chapter body parse error")
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Bulk delete chapter validation failed")
	}

	if err := c.chapterUsecase.BulkDelete(ctx.UserContext(), req.IDs); err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, fiber.Map{"message": "chapters deleted successfully"})
}

// GetByBookID handles getting chapters by book ID
// GET /api/v1/chapters/book/:bookId
func (c *ChapterController) GetByBookID(ctx *fiber.Ctx) error {
	bookID, err := strconv.Atoi(ctx.Params("bookId"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid book ID", err, nil, "")
	}

	result, err := c.chapterUsecase.GetByBookID(ctx.UserContext(), uint(bookID))
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	out := make([]dto.ListChapterResponse, 0, len(result.Data))
	for _, chapter := range result.Data {
		out = append(out, c.toListChapterResponse(&chapter))
	}

	return response.SendPaginated(ctx, out, result.Page, result.Limit, result.Total, result.TotalPages, len(result.Data))
}
