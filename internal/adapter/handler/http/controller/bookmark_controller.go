package controller

import (
	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/adapter/handler/http/middleware"
	"ishari-backend/internal/adapter/handler/http/response"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/pkg/errors"
	"ishari-backend/pkg/logger"
	"ishari-backend/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

type BookmarkController struct {
	bookmarkUsecase portuc.BookmarkUsecase
	validator       validation.Validator
	log             logger.Logger
}

func NewBookmarkController(bookmarkUsecase portuc.BookmarkUsecase, validator validation.Validator, log logger.Logger) *BookmarkController {
	return &BookmarkController{
		bookmarkUsecase: bookmarkUsecase,
		validator:       validator,
		log:             log,
	}
}

// Create Bookmark
func (c *BookmarkController) Create(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromContext(ctx)
	if user == nil {
		return errors.Unauthorized("unauthorized: missing user")
	}
	userID := user.UserID

	var req dto.CreateBookmarkRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.log.Error("failed to parse request body", "error", err)
		return errors.BadRequest("invalid request body")
	}

	if errs := c.validator.Struct(req); errs != nil {
		return errors.BadRequest("validation failed")
	}

	input := portuc.CreateBookmarkInput{
		UserID:  userID,
		VerseID: req.VerseID,
		Note:    req.Note,
	}

	bookmark, err := c.bookmarkUsecase.Create(ctx.UserContext(), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	resp := dto.BookmarkResponse{
		ID:        bookmark.ID,
		UserID:    bookmark.UserID,
		VerseID:   bookmark.VerseID,
		Note:      bookmark.Note,
		CreatedAt: bookmark.CreatedAt,
	}

	return response.SendCreated(ctx, "bookmark created successfully", resp)
}

// Get Bookmark By ID
func (c *BookmarkController) GetByID(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromContext(ctx)
	if user == nil {
		return errors.Unauthorized("unauthorized: missing user")
	}
	userID := user.UserID

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errors.BadRequest("invalid bookmark id")
	}

	bookmark, err := c.bookmarkUsecase.GetByID(ctx.UserContext(), uint(id), userID)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	resp := dto.BookmarkResponse{
		ID:        bookmark.ID,
		UserID:    bookmark.UserID,
		VerseID:   bookmark.VerseID,
		Note:      bookmark.Note,
		CreatedAt: bookmark.CreatedAt,
	}

	return response.SendOK(ctx, resp)
}

// List Bookmarks By User ID
func (c *BookmarkController) List(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromContext(ctx)
	if user == nil {
		return errors.Unauthorized("unauthorized: missing user")
	}
	userID := user.UserID

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)
	sort := ctx.Query("sort", "created_at desc")

	input := portuc.ListBookmarkInput{
		UserID: userID,
		Page:   page,
		Limit:  limit,
		Sort:   sort,
	}

	paginatedResult, err := c.bookmarkUsecase.ListByUserID(ctx.UserContext(), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	var respItems []dto.ListBookmarkResponse
	for _, b := range paginatedResult.Data {
		respItems = append(respItems, dto.ListBookmarkResponse{
			ID:        b.ID,
			UserID:    b.UserID,
			VerseID:   b.VerseID,
			Note:      b.Note,
			CreatedAt: b.CreatedAt,
		})
	}

	// For empty arrays, ensure we return [] instead of null in JSON
	if respItems == nil {
		respItems = make([]dto.ListBookmarkResponse, 0)
	}

	return response.SendPaginated(ctx, respItems, paginatedResult.Page, paginatedResult.Limit, paginatedResult.Total, paginatedResult.TotalPages, len(respItems))
}

// List Bookmarks By Specific User ID
func (c *BookmarkController) ListBySpecificUser(ctx *fiber.Ctx) error {
	// Parse user ID from path parameter
	targetUserID, err := ctx.ParamsInt("userId")
	if err != nil {
		return errors.BadRequest("invalid user id")
	}

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)
	sort := ctx.Query("sort", "created_at desc")

	input := portuc.ListBookmarkInput{
		UserID: uint(targetUserID),
		Page:   page,
		Limit:  limit,
		Sort:   sort,
	}

	paginatedResult, err := c.bookmarkUsecase.ListByUserID(ctx.UserContext(), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	var respItems []dto.ListBookmarkResponse
	for _, b := range paginatedResult.Data {
		respItems = append(respItems, dto.ListBookmarkResponse{
			ID:        b.ID,
			UserID:    b.UserID,
			VerseID:   b.VerseID,
			Note:      b.Note,
			CreatedAt: b.CreatedAt,
		})
	}

	// For empty arrays, ensure we return [] instead of null in JSON
	if respItems == nil {
		respItems = make([]dto.ListBookmarkResponse, 0)
	}

	return response.SendPaginated(ctx, respItems, paginatedResult.Page, paginatedResult.Limit, paginatedResult.Total, paginatedResult.TotalPages, len(respItems))
}

// Update Bookmark
func (c *BookmarkController) Update(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromContext(ctx)
	if user == nil {
		return errors.Unauthorized("unauthorized: missing user")
	}
	userID := user.UserID

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errors.BadRequest("invalid bookmark id")
	}

	var req dto.UpdateBookmarkRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.log.Error("failed to parse request body", "error", err)
		return errors.BadRequest("invalid request body")
	}

	if errs := c.validator.Struct(req); errs != nil {
		return errors.BadRequest("validation failed")
	}

	input := portuc.UpdateBookmarkInput{
		UserID: userID,
		Note:   req.Note,
	}

	bookmark, err := c.bookmarkUsecase.Update(ctx.UserContext(), uint(id), input)
	if err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	resp := dto.BookmarkResponse{
		ID:        bookmark.ID,
		UserID:    bookmark.UserID,
		VerseID:   bookmark.VerseID,
		Note:      bookmark.Note,
		CreatedAt: bookmark.CreatedAt,
	}

	return response.SendOK(ctx, resp)
}

// Delete Bookmark
func (c *BookmarkController) Delete(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromContext(ctx)
	if user == nil {
		return errors.Unauthorized("unauthorized: missing user")
	}
	userID := user.UserID

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errors.BadRequest("invalid bookmark id")
	}

	if err := c.bookmarkUsecase.Delete(ctx.UserContext(), uint(id), userID); err != nil {
		return response.SendDomainError(ctx, err, c.log)
	}

	return response.SendOK(ctx, fiber.Map{
		"message": "bookmark deleted successfully",
	})
}
