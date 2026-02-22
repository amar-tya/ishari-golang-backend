package chapter

import (
	"context"

	"ishari-backend/internal/core/domain"
	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/logger"
	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
)

type chapterUsecase struct {
	chapterRepo repository.ChapterRepository
	bookRepo    repository.BookRepository
	log         logger.Logger
}

func NewChapterUsecase(chapterRepo repository.ChapterRepository, bookRepo repository.BookRepository, log logger.Logger) portuc.ChapterUsecase {
	return &chapterUsecase{
		chapterRepo: chapterRepo,
		bookRepo:    bookRepo,
		log:         log,
	}
}

// Create creates a new chapter
func (u *chapterUsecase) Create(ctx context.Context, input portuc.CreateChapterInput) (*entity.Chapter, error) {

	// Validate input
	if err := u.validateCreateChapter(ctx, input); err != nil {
		return nil, err
	}

	// create chapter entity
	chapter := &entity.Chapter{
		BookID:        input.BookID,
		ChapterNumber: input.ChapterNumber,
		Title:         input.Title,
		Category:      input.Category,
		Description:   input.Description,
		TotalVerses:   input.TotalVerses,
	}

	// Persist chapter
	if err := u.chapterRepo.CreateChapter(ctx, chapter); err != nil {
		return nil, err
	}

	return chapter, nil
}

// validateCreateChapter validates create chapter input
func (u *chapterUsecase) validateCreateChapter(ctx context.Context, input portuc.CreateChapterInput) error {
	if err := u.validateBookId(ctx, input.BookID); err != nil {
		return err
	}
	if err := u.validateChapterNumber(input.ChapterNumber); err != nil {
		return err
	}
	if err := u.validateTitle(input.Title); err != nil {
		return err
	}
	if err := u.validateCategory(input.Category); err != nil {
		return err
	}
	if err := u.validateTotalVerses(input.TotalVerses); err != nil {
		return err
	}
	return nil
}

// validateBookId checks if the book ID exists
func (u *chapterUsecase) validateBookId(ctx context.Context, bookId uint) error {
	book, err := u.bookRepo.GetById(ctx, int64(bookId))
	if err != nil {
		u.log.Error("failed to get book", "error", err, "book_id", bookId)
		return domain.NewInternalError("failed to validate book", err)
	}
	if book == nil {
		return ErrBookNotFound
	}
	return nil
}

// validateChapterNumber checks if the chapter number is valid
func (u *chapterUsecase) validateChapterNumber(chapterNumber uint) error {
	if chapterNumber < 1 {
		return ErrInvalidChapterNumber
	}
	return nil
}

// validateTitle checks if the title is valid
func (u *chapterUsecase) validateTitle(title string) error {
	if title == "" {
		return ErrInvalidTitle
	}
	return nil
}

// validateCategory checks if the category is valid
func (u *chapterUsecase) validateCategory(category string) error {
	if category == "" {
		return ErrInvalidCategory
	}
	return nil
}

// validateTotalVerses checks if the total verses is valid
func (u *chapterUsecase) validateTotalVerses(totalVerses uint) error {
	if totalVerses < 1 {
		return ErrInvalidTotalVerses
	}
	return nil
}

// List returns paginated chapters with optional search
func (u *chapterUsecase) List(ctx context.Context, params portuc.ListChapterInput) (*portuc.PaginatedResult[entity.Chapter], error) {
	// apply param default
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit < 0 {
		params.Limit = 20
	}

	offset := (params.Page - 1) * params.Limit

	chapters, total, err := u.chapterRepo.ListChapters(ctx, offset, params.Limit, params.Search, params.BookID, params.Title, params.Category)
	if err != nil {
		u.log.Error("failed to list chapters", "error", err)
		return nil, domain.NewInternalError("failed to list chapters", err)
	}

	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}

	return &portuc.PaginatedResult[entity.Chapter]{
		Data:       chapters,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil

}

// GetByID returns a chapter by ID
func (u *chapterUsecase) GetByID(ctx context.Context, id uint) (*entity.Chapter, error) {
	chapter, err := u.chapterRepo.GetChapterByID(ctx, id)
	if err != nil {
		u.log.Error("failed to get chapter by ID", "error", err, "chapter_id", id)
		return nil, ErrChapterNotFound
	}
	return chapter, nil
}

// GetByBookID returns all chapters for a given book ID with pagination
func (u *chapterUsecase) GetByBookID(ctx context.Context, bookID uint) (*portuc.PaginatedResult[entity.Chapter], error) {
	chapters, total, err := u.chapterRepo.GetChaptersByBookID(ctx, bookID)
	if err != nil {
		u.log.Error("failed to get chapters by book ID", "error", err, "book_id", bookID)
		return nil, domain.NewInternalError("failed to get chapters by book ID", err)
	}

	return &portuc.PaginatedResult[entity.Chapter]{
		Data:       chapters,
		Total:      total,
		Page:       1,
		Limit:      int(total),
		TotalPages: 1,
	}, nil
}

// Update modifies chapter information
func (u *chapterUsecase) Update(ctx context.Context, id uint, input portuc.UpdateChapterInput) (*entity.Chapter, error) {
	chapter, err := u.chapterRepo.GetChapterByID(ctx, id)
	if err != nil {
		u.log.Error("failed to get chapter for update", "error", err, "chapter_id", id)
		return nil, ErrChapterNotFound
	}

	// update fields if provided
	if input.BookID != nil {
		if err := u.validateBookId(ctx, *input.BookID); err != nil {
			return nil, err
		}
		chapter.BookID = *input.BookID
	}
	if input.ChapterNumber != nil {
		if err := u.validateChapterNumber(*input.ChapterNumber); err != nil {
			return nil, err
		}
		chapter.ChapterNumber = *input.ChapterNumber
	}
	if input.Title != nil {
		if err := u.validateTitle(*input.Title); err != nil {
			return nil, err
		}
		chapter.Title = *input.Title
	}

	if input.Category != nil {
		if err := u.validateCategory(*input.Category); err != nil {
			return nil, err
		}
		chapter.Category = *input.Category
	}

	if input.TotalVerses != nil {
		if err := u.validateTotalVerses(*input.TotalVerses); err != nil {
			return nil, err
		}
		chapter.TotalVerses = *input.TotalVerses
	}

	if err := u.chapterRepo.UpdateChapter(ctx, chapter); err != nil {
		u.log.Error("failed to update chapter", "error", err, "chapter_id", id)
		return nil, domain.NewInternalError("failed to update chapter", err)
	}

	return chapter, nil
}

// Delete removes a chapter by ID
func (u *chapterUsecase) Delete(ctx context.Context, id uint) error {
	// Check if chapter exists
	_, err := u.chapterRepo.GetChapterByID(ctx, id)
	if err != nil {
		u.log.Error("chapter not found for deletion", "error", err, "chapter_id", id)
		return ErrChapterNotFound
	}

	// Delete the chapter
	if err := u.chapterRepo.DeleteChapter(ctx, id); err != nil {
		u.log.Error("failed to delete chapter", "error", err, "chapter_id", id)
		return domain.NewInternalError("failed to delete chapter", err)
	}
	return nil
}

// BulkDelete removes multiple chapters by IDs
func (u *chapterUsecase) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	if err := u.chapterRepo.DeleteChapters(ctx, ids); err != nil {
		u.log.Error("failed to bulk delete chapters", "error", err, "chapter_ids", ids)
		return domain.NewInternalError("failed to bulk delete chapters", err)
	}

	return nil
}
