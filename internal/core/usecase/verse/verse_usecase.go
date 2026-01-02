package verse

import (
	"context"
	"ishari-backend/internal/core/domain"
	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/logger"
	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
)

type verseUsecase struct {
	verseRepo   repository.VerseRepository
	chapterRepo repository.ChapterRepository
	log         logger.Logger
}

func NewVerseUsecase(verRepo repository.VerseRepository, chapterRepo repository.ChapterRepository, log logger.Logger) portuc.VerseUseCase {
	return &verseUsecase{
		verseRepo:   verRepo,
		chapterRepo: chapterRepo,
		log:         log,
	}
}

// Create creates a new verse
func (u *verseUsecase) Create(ctx context.Context, input portuc.CreateVerseInput) (*entity.Verse, error) {

	// Validate input
	if err := u.validateCreateVerse(ctx, input); err != nil {
		return nil, err
	}

	// create verse entity
	verse := &entity.Verse{
		ChapterID:       input.ChapterID,
		VerseNumber:     input.VerseNumber,
		ArabicText:      input.ArabicText,
		Transliteration: input.Transliteration,
	}

	// Persist verse
	if err := u.verseRepo.Create(ctx, verse); err != nil {
		return nil, err
	}

	return verse, nil
}

func (u *verseUsecase) validateCreateVerse(ctx context.Context, input portuc.CreateVerseInput) error {
	if input.ChapterID == 0 {
		return ErrChapterNotFound
	}
	if err := u.validateChapterId(ctx, input.ChapterID); err != nil {
		return err
	}
	if input.VerseNumber == 0 {
		return ErrInvalidVerseNumber
	}
	if err := u.validateVerseNumber(input.VerseNumber); err != nil {
		return err
	}
	if input.ArabicText == "" {
		return ErrInvalidVerseText
	}
	if err := u.validateVerseText(input.ArabicText); err != nil {
		return err
	}
	return nil
}

func (u *verseUsecase) validateChapterId(ctx context.Context, chapterId uint) error {
	chapter, err := u.chapterRepo.GetChapterByID(ctx, chapterId)
	if err != nil {
		u.log.Error("failed to get chapter", "error", err, "chapter_id", chapterId)
		return domain.NewInternalError("failed to validate chapter", err)
	}
	if chapter == nil {
		return ErrChapterNotFound
	}
	return nil
}

func (u *verseUsecase) validateVerseNumber(verseNumber uint) error {
	if verseNumber == 0 {
		return ErrInvalidVerseNumber
	}
	return nil
}

func (u *verseUsecase) validateVerseText(verseText string) error {
	if verseText == "" {
		return ErrInvalidVerseText
	}
	return nil
}

func (u *verseUsecase) List(ctx context.Context, params portuc.ListParams) (*portuc.PaginatedResult[entity.Verse], error) {

	// apply param default
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}

	offset := (params.Page - 1) * params.Limit

	verses, total, err := u.verseRepo.List(ctx, offset, params.Limit, params.Search)
	if err != nil {
		u.log.Error("failed to list verses", "error", err)
		return nil, domain.NewInternalError("failed to list verses", err)
	}

	totalPages := int(total / uint(params.Limit))
	if uint(total%uint(params.Limit)) > 0 {
		totalPages++
	}
	return &portuc.PaginatedResult[entity.Verse]{
		Data:       verses,
		Total:      int64(total),
		Page:       int(params.Page),
		Limit:      int(params.Limit),
		TotalPages: totalPages,
	}, nil
}

func (u *verseUsecase) GetById(ctx context.Context, id uint) (*entity.Verse, error) {
	verse, err := u.verseRepo.GetById(ctx, id)
	if err != nil {
		u.log.Error("failed to get verse", "error", err, "verse_id", id)
		return nil, domain.NewInternalError("failed to get verse", err)
	}
	if verse == nil {
		return nil, ErrVerseNotFound
	}
	return verse, nil
}

func (u *verseUsecase) Update(ctx context.Context, id uint, input portuc.UpdateVerseInput) (*entity.Verse, error) {
	verse, err := u.verseRepo.GetById(ctx, id)
	if err != nil {
		u.log.Error("failed to get verse for update", "error", err, "verse_id", id)
		return nil, ErrVerseNotFound
	}

	// update fields if provided
	if input.ChapterID != nil {
		if err := u.validateChapterId(ctx, *input.ChapterID); err != nil {
			return nil, err
		}
		verse.ChapterID = *input.ChapterID
	}
	if input.VerseNumber != nil {
		if err := u.validateVerseNumber(*input.VerseNumber); err != nil {
			return nil, err
		}
		verse.VerseNumber = *input.VerseNumber
	}
	if input.ArabicText != nil {
		if err := u.validateVerseText(*input.ArabicText); err != nil {
			return nil, err
		}
		verse.ArabicText = *input.ArabicText
	}
	if input.Transliteration != nil {
		verse.Transliteration = input.Transliteration
	}

	// update verse
	if err := u.verseRepo.Update(ctx, verse); err != nil {
		u.log.Error("failed to update verse", "error", err, "verse_id", id)
		return nil, domain.NewInternalError("failed to update verse", err)
	}

	return verse, nil
}

func (u *verseUsecase) Delete(ctx context.Context, id uint) error {
	_, err := u.verseRepo.GetById(ctx, id)
	if err != nil {
		u.log.Error("failed to get verse for delete", "error", err, "verse_id", id)
		return ErrVerseNotFound
	}

	if err := u.verseRepo.Delete(ctx, id); err != nil {
		u.log.Error("failed to delete verse", "error", err, "verse_id", id)
		return domain.NewInternalError("failed to delete verse", err)
	}
	return nil
}
