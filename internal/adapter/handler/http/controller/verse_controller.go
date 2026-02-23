package controller

import (
	"encoding/json"
	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/adapter/handler/http/response"
	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/logger"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/pkg/validation"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type VerseController struct {
	verseUsecase portuc.VerseUseCase
	validate     validation.Validator
	log          logger.Logger
}

func NewVerseController(verseUsecase portuc.VerseUseCase, validate validation.Validator, log logger.Logger) *VerseController {
	return &VerseController{
		verseUsecase: verseUsecase,
		validate:     validate,
		log:          log,
	}
}

// Create handles creating a new verse
// POST /api/verses
func (c *VerseController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateVerseRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.SendParseError(ctx, err, c.log, "Create verse body parse error")
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Create verse validation failed")
	}

	input := portuc.CreateVerseInput{
		ChapterID:       req.ChapterID,
		VerseNumber:     req.VerseNumber,
		ArabicText:      req.ArabicText,
		Transliteration: req.Transliteration,
	}

	verse, err := c.verseUsecase.Create(ctx.UserContext(), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendCreated(ctx, "verse created successfully", c.toListVerseResponse(verse))
}

func (c *VerseController) toListVerseResponse(verse *entity.Verse) dto.ListVerseResponse {
	var chapterResp *dto.ListChapterResponse
	if verse.Chapter != nil {
		chapterResp = &dto.ListChapterResponse{
			ID:            verse.Chapter.ID,
			BookID:        verse.Chapter.BookID,
			ChapterNumber: verse.Chapter.ChapterNumber,
			Title:         verse.Chapter.Title,
			Category:      verse.Chapter.Category,
			Description:   verse.Chapter.Description,
			TotalVerses:   verse.Chapter.TotalVerses,
			CreatedAt:     verse.Chapter.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:     verse.Chapter.UpdatedAt.UTC().Format(time.RFC3339),
		}

		if verse.Chapter.Book != nil {
			chapterResp.Book = &dto.BookResponse{
				ID:            verse.Chapter.Book.ID,
				Title:         verse.Chapter.Book.Title,
				Author:        verse.Chapter.Book.Author,
				Description:   verse.Chapter.Book.Description,
				PublishedYear: verse.Chapter.Book.PublishedYear,
				CoverImageURL: verse.Chapter.Book.CoverImageURL,
				CreatedAt:     verse.Chapter.Book.CreatedAt.UTC().Format(time.RFC3339),
				UpdatedAt:     verse.Chapter.Book.UpdatedAt.UTC().Format(time.RFC3339),
			}
		}
	}

	return dto.ListVerseResponse{
		ID:              verse.ID,
		ChapterID:       verse.ChapterID,
		Chapter:         chapterResp,
		VerseNumber:     verse.VerseNumber,
		ArabicText:      verse.ArabicText,
		Transliteration: verse.Transliteration,
		CreatedAt:       verse.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       verse.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

// List handles listing verses
// GET /api/verses
func (c *VerseController) List(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))
	search := ctx.Query("search", "")

	var chapterID *uint
	if cidStr := ctx.Query("chapter_id"); cidStr != "" {
		if cid, err := strconv.Atoi(cidStr); err == nil && cid > 0 {
			u := uint(cid)
			chapterID = &u
		}
	}

	arabicText := ctx.Query("arabic_text", "")
	transliteration := ctx.Query("transliteration", "")

	params := portuc.ListParams{
		Page:            uint(page),
		Limit:           uint(limit),
		Search:          search,
		ChapterID:       chapterID,
		ArabicText:      arabicText,
		Transliteration: transliteration,
	}

	result, err := c.verseUsecase.List(ctx.UserContext(), params)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	var totalPages int
	if limit > 0 {
		totalPages = int(math.Ceil(float64(result.Total) / float64(limit)))
	}

	out := make([]dto.ListVerseResponse, 0, len(result.Data))
	for _, verse := range result.Data {
		out = append(out, c.toListVerseResponse(&verse))
	}

	return response.SendPaginated(ctx, out, page, limit, result.Total, totalPages, len(result.Data))
}

// GetByID handles getting a verse by ID
// GET /api/verses/:id
func (c *VerseController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid verse ID", err, nil, "")
	}

	verse, err := c.verseUsecase.GetById(ctx.UserContext(), uint(id))
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, c.toListVerseResponse(verse))
}

// Update handles updating a verse
// PUT /api/verses/:id
func (c *VerseController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid verse ID", err, nil, "")
	}

	var req dto.UpdateVerseRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.SendParseError(ctx, err, c.log, "Update verse body parse error")
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Update verse validation failed")
	}

	input := portuc.UpdateVerseInput{
		ChapterID:       req.ChapterID,
		VerseNumber:     req.VerseNumber,
		ArabicText:      req.ArabicText,
		Transliteration: req.Transliteration,
	}

	verse, err := c.verseUsecase.Update(ctx.UserContext(), uint(id), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, c.toListVerseResponse(verse))
}

// Delete handles deleting a verse
// DELETE /api/verses/:id
func (c *VerseController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid verse ID", err, nil, "")
	}

	err = c.verseUsecase.Delete(ctx.UserContext(), uint(id))
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, "verse deleted successfully")
}

// BulkDelete handles bulk deleting verses
// POST /api/verses/bulk-delete
func (c *VerseController) BulkDelete(ctx *fiber.Ctx) error {
	var req dto.BulkDeleteVerseRequest

	// Try to parse body first
	if err := ctx.BodyParser(&req); err != nil {
		c.log.Error("BodyParser failed", err)
	}

	// If IDs are still empty but body exists, try manual unmarshal
	// This handles cases where BodyParser returns nil but didn't parse (e.g. missing Content-Type)
	if len(req.IDs) == 0 && len(ctx.Body()) > 0 {
		if jsonErr := json.Unmarshal(ctx.Body(), &req); jsonErr != nil {
			c.log.Error("Manual JSON unmarshal failed", jsonErr)
		}
	}

	// If no IDs found in body, try query params
	if len(req.IDs) == 0 {
		if err := ctx.QueryParser(&req); err != nil {
			return response.SendParseError(ctx, err, c.log, "Bulk delete verse query parse error")
		}
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Bulk delete verse validation failed")
	}

	if err := c.verseUsecase.BulkDelete(ctx.UserContext(), req.IDs); err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, fiber.Map{"message": "verses deleted successfully"})
}
