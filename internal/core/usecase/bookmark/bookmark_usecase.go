package bookmark

import (
	"context"

	"ishari-backend/internal/core/domain"
	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/logger"
	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
)

var (
	ErrBookmarkNotFound      = domain.NewNotFoundError("bookmark not found", nil)
	ErrBookmarkAlreadyExists = domain.NewConflictError("bookmark for this verse already exists", nil)
	ErrInvalidUserId         = domain.NewInvalidInputError("invalid user ID", nil)
	ErrInvalidVerseId        = domain.NewInvalidInputError("invalid verse ID", nil)
	ErrForbidden             = domain.NewUnauthorizedError("you don't have permission to modify this bookmark", nil)
)

type bookmarkUsecase struct {
	bookmarkRepo repository.BookmarkRepository
	verseRepo    repository.VerseRepository
	log          logger.Logger
}

func NewBookmarkUsecase(bookmarkRepo repository.BookmarkRepository, verseRepo repository.VerseRepository, log logger.Logger) portuc.BookmarkUsecase {
	return &bookmarkUsecase{
		bookmarkRepo: bookmarkRepo,
		verseRepo:    verseRepo,
		log:          log,
	}
}

func (u *bookmarkUsecase) Create(ctx context.Context, input portuc.CreateBookmarkInput) (*entity.Bookmark, error) {
	if input.UserID == 0 {
		return nil, ErrInvalidUserId
	}
	if input.VerseID == 0 {
		return nil, ErrInvalidVerseId
	}

	// Validate Verse exists
	_, err := u.verseRepo.GetById(ctx, input.VerseID)
	if err != nil {
		u.log.Error("failed to get verse", "error", err, "verse_id", input.VerseID)
		// Assume if there is an error it might be verse not found (could refine with specific VerseNotFound error later)
		return nil, domain.NewInvalidInputError("verse not found", err)
	}

	// Check if already bookmarked
	existing, err := u.bookmarkRepo.GetBookmarkByUserIDAndVerseID(ctx, input.UserID, input.VerseID)
	if err != nil {
		u.log.Error("failed to check existing bookmark", "error", err)
		return nil, domain.NewInternalError("failed to check existing bookmark", err)
	}
	if existing != nil {
		return nil, ErrBookmarkAlreadyExists
	}

	bookmark := &entity.Bookmark{
		UserID:  input.UserID,
		VerseID: input.VerseID,
		Note:    input.Note,
	}

	if err := u.bookmarkRepo.CreateBookmark(ctx, bookmark); err != nil {
		u.log.Error("failed to create bookmark", "error", err)
		return nil, domain.NewInternalError("failed to create bookmark", err)
	}

	return bookmark, nil
}

func (u *bookmarkUsecase) GetByID(ctx context.Context, id uint, userID uint) (*entity.Bookmark, error) {
	bookmark, err := u.bookmarkRepo.GetBookmarkByID(ctx, id)
	if err != nil {
		u.log.Error("failed to get bookmark", "error", err, "bookmark_id", id)
		return nil, ErrBookmarkNotFound
	}

	if bookmark.UserID != userID {
		return nil, ErrForbidden
	}

	return bookmark, nil
}

func (u *bookmarkUsecase) ListByUserID(ctx context.Context, input portuc.ListBookmarkInput) (*portuc.PaginatedResult[entity.Bookmark], error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}

	offset := (input.Page - 1) * input.Limit

	bookmarks, total, err := u.bookmarkRepo.ListBookmarksByUserID(ctx, input.UserID, offset, input.Limit, input.Sort)
	if err != nil {
		u.log.Error("failed to list bookmarks", "error", err, "user_id", input.UserID)
		return nil, domain.NewInternalError("failed to list bookmarks", err)
	}

	totalPages := int(total) / input.Limit
	if int(total)%input.Limit > 0 {
		totalPages++
	}

	return &portuc.PaginatedResult[entity.Bookmark]{
		Data:       bookmarks,
		Total:      total,
		Page:       input.Page,
		Limit:      input.Limit,
		TotalPages: totalPages,
	}, nil
}

func (u *bookmarkUsecase) Update(ctx context.Context, id uint, input portuc.UpdateBookmarkInput) (*entity.Bookmark, error) {
	bookmark, err := u.bookmarkRepo.GetBookmarkByID(ctx, id)
	if err != nil {
		u.log.Error("failed to get bookmark for update", "error", err, "bookmark_id", id)
		return nil, ErrBookmarkNotFound
	}

	// verify ownership
	if bookmark.UserID != input.UserID {
		return nil, ErrForbidden
	}

	// Update logic: Only Note is allowed to be updated.
	bookmark.Note = input.Note

	if err := u.bookmarkRepo.UpdateBookmark(ctx, bookmark); err != nil {
		u.log.Error("failed to update bookmark", "error", err, "bookmark_id", id)
		return nil, domain.NewInternalError("failed to update bookmark", err)
	}

	return bookmark, nil
}

func (u *bookmarkUsecase) Delete(ctx context.Context, id uint, userID uint) error {
	bookmark, err := u.bookmarkRepo.GetBookmarkByID(ctx, id)
	if err != nil {
		u.log.Error("failed to get bookmark for deletion", "error", err, "bookmark_id", id)
		return ErrBookmarkNotFound
	}

	// verify ownership
	if bookmark.UserID != userID {
		return ErrForbidden
	}

	if err := u.bookmarkRepo.DeleteBookmark(ctx, id); err != nil {
		u.log.Error("failed to delete bookmark", "error", err, "bookmark_id", id)
		return domain.NewInternalError("failed to delete bookmark", err)
	}

	return nil
}
