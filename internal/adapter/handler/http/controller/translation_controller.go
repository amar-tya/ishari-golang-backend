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

type TranslationController struct {
	translationUsecase portuc.TranslationUseCase
	validate           validation.Validator
	log                logger.Logger
}

func NewTranslationController(translationUsecase portuc.TranslationUseCase, validate validation.Validator, log logger.Logger) *TranslationController {
	return &TranslationController{
		translationUsecase: translationUsecase,
		validate:           validate,
		log:                log,
	}
}

// Create handles creating a new translation
// POST /api/translations
func (c *TranslationController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateTranslationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.SendParseError(ctx, err, c.log, "Create translation body parse error")
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Create translation validation faildes")
	}

	input := portuc.CreateTranslationInput{
		VerseID:         req.VerseID,
		LanguageCode:    req.LanguageCode,
		TranslationText: req.TranslationText,
		TranslatorName:  &req.TranslatorName,
	}

	translation, err := c.translationUsecase.Create(ctx.UserContext(), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendCreated(ctx, "translation created successfully", c.toListTranslationResponse(translation))
}

// List handles listing translations
// GET /api/translations
func (c *TranslationController) List(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))
	search := ctx.Query("search", "")

	verseIDStr := ctx.Query("verse_id", ctx.Query("verseId", "0"))
	verseID, _ := strconv.Atoi(verseIDStr)
	translatorName := ctx.Query("translator_name", ctx.Query("translatorName", ""))
	languageCode := ctx.Query("language_code", ctx.Query("languageCode", ""))

	params := portuc.TranslationListParams{
		Page:           uint(page),
		Limit:          uint(limit),
		Search:         search,
		VerseID:        uint(verseID),
		TranslatorName: translatorName,
		LanguageCode:   languageCode,
	}

	result, err := c.translationUsecase.List(ctx.UserContext(), params)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	var totalPages int
	if limit > 0 {
		totalPages = int(math.Ceil(float64(result.Total) / float64(limit)))
	}

	out := make([]dto.ListTranslationResponse, 0, len(result.Data))
	for i := range result.Data {
		out = append(out, c.toListTranslationResponse(&result.Data[i]))
	}

	return response.SendPaginated(ctx, out, page, limit, result.Total, totalPages, len(result.Data))
}

// GetByID handles getting a translation by ID
// GET /api/translations/:id
func (c *TranslationController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid translation ID", err, nil, "")
	}

	translation, err := c.translationUsecase.GetById(ctx.UserContext(), uint(id))
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, c.toListTranslationResponse(translation))
}

// GetByVerseID handles getting translations by verse ID
// GET /api/translations/verse/:verse_id
func (c *TranslationController) GetByVerseID(ctx *fiber.Ctx) error {
	verseID, err := strconv.Atoi(ctx.Params("verse_id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid verse ID", err, nil, "")
	}

	translations, err := c.translationUsecase.GetByVerseId(ctx.UserContext(), uint(verseID))
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	out := make([]dto.ListTranslationResponse, 0, len(translations))
	for i := range translations {
		out = append(out, c.toListTranslationResponse(&translations[i]))
	}

	return response.SendOK(ctx, out)
}

// Update handles updating a translation
// PUT /api/translations/:id
func (c *TranslationController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid translation ID", err, nil, "")
	}

	var req dto.UpdateTranslationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.SendParseError(ctx, err, c.log, "Update translation body parse error")
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Update translation validation failed")
	}

	input := portuc.UpdateTranslationInput{
		VerseID:         req.VerseID,
		LanguageCode:    req.LanguageCode,
		TranslationText: req.TranslationText,
		TranslatorName:  req.TranslatorName,
	}

	translation, err := c.translationUsecase.Update(ctx.UserContext(), uint(id), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, c.toListTranslationResponse(translation))
}

// Delete handles deleting a translation
// DELETE /api/translations/:id
func (c *TranslationController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return response.SendBadRequest(ctx, "invalid translation ID", err, nil, "")
	}

	err = c.translationUsecase.Delete(ctx.UserContext(), uint(id))
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, "translation deleted successfully")
}

// BulkDelete handles bulk deleting translations
// POST /api/translations/bulk-delete
func (c *TranslationController) BulkDelete(ctx *fiber.Ctx) error {
	var req dto.BulkDeleteTranslationRequest

	// Try to parse body first
	if err := ctx.BodyParser(&req); err != nil {
		c.log.Error("BodyParser failed", err)
	}

	// If IDs are still empty but body exists, try manual unmarshal
	if len(req.IDs) == 0 && len(ctx.Body()) > 0 {
		if jsonErr := json.Unmarshal(ctx.Body(), &req); jsonErr != nil {
			c.log.Error("Manual JSON unmarshal failed", jsonErr)
		}
	}

	// If no IDs found in body, try query params
	if len(req.IDs) == 0 {
		if err := ctx.QueryParser(&req); err != nil {
			return response.SendParseError(ctx, err, c.log, "Bulk delete translation query parse error")
		}
	}

	if err := c.validate.Struct(req); err != nil {
		return response.SendValidationError(ctx, err, c.log, "Bulk delete translation validation failed")
	}

	if err := c.translationUsecase.BulkDelete(ctx.UserContext(), req.IDs); err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, fiber.Map{"message": "translations deleted successfully"})
}

func (c *TranslationController) toListTranslationResponse(translation *entity.Translation) dto.ListTranslationResponse {
	out := dto.ListTranslationResponse{
		ID:              translation.ID,
		VerseID:         translation.VerseID,
		LanguageCode:    translation.LanguageCode,
		TranslationText: translation.TranslationText,
		TranslatorName:  translation.TranslatorName,
		CreatedAt:       translation.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       translation.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if translation.Verse != nil {
		out.Verse = &dto.VerseDropdownItem{
			ID:         translation.Verse.ID,
			ArabicText: translation.Verse.ArabicText,
		}
	}
	return out
}

// GetDropdown handles getting dropdown filter data for translations
// GET /api/translations/dropdown
func (c *TranslationController) GetDropdown(ctx *fiber.Ctx) error {
	data, err := c.translationUsecase.GetDropdownData(ctx.UserContext())
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	// Build verse dropdown items
	verseItems := make([]dto.VerseDropdownItem, 0, len(data.Verses))
	for _, v := range data.Verses {
		verseItems = append(verseItems, dto.VerseDropdownItem{
			ID:         v.ID,
			ArabicText: v.ArabicText,
		})
	}

	return response.SendOK(ctx, dto.TranslationDropdownResponse{
		Verses:          verseItems,
		TranslatorNames: data.TranslatorNames,
		LanguageCodes:   data.LanguageCodes,
	})
}
