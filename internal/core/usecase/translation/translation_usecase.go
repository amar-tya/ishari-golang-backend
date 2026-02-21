package translation

import (
	"context"
	"ishari-backend/internal/core/domain"
	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/logger"
	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
)

type translationUsecase struct {
	translationRepository repository.TranslationRepository
	verseRepository       repository.VerseRepository
	log                   logger.Logger
}

func NewTranslationUsecase(transRepo repository.TranslationRepository, verRepo repository.VerseRepository, log logger.Logger) portuc.TranslationUseCase {
	return &translationUsecase{
		translationRepository: transRepo,
		verseRepository:       verRepo,
		log:                   log,
	}
}

// Create a new translation
func (u *translationUsecase) Create(ctx context.Context, input portuc.CreateTranslationInput) (*entity.Translation, error) {
	// validate input
	if err := u.validateCreateTranslation(ctx, input); err != nil {
		return nil, err
	}

	// create translation entity
	translation := &entity.Translation{
		VerseID:         input.VerseID,
		LanguageCode:    input.LanguageCode,
		TranslationText: input.TranslationText,
		TranslatorName:  input.TranslatorName,
	}

	// Persist translation
	if err := u.translationRepository.Create(ctx, translation); err != nil {
		return nil, err
	}

	// Reload with preload to get Verse data for the response
	reloaded, err := u.translationRepository.GetById(ctx, translation.ID)
	if err != nil {
		// Log error but return the translation anyway (Verse will be nil)
		u.log.Error("failed to reload translation after creation", "error", err, "id", translation.ID)
		return translation, nil
	}

	return reloaded, nil
}

func (u *translationUsecase) validateCreateTranslation(ctx context.Context, input portuc.CreateTranslationInput) error {
	if input.VerseID == 0 {
		return ErrTranslationNotFound
	}
	if err := u.validateVerseId(ctx, input.VerseID); err != nil {
		return err
	}
	if input.LanguageCode == "" {
		return ErrInvalidLanguageCode
	}
	if input.TranslationText == "" {
		return ErrInvalidTranslationText
	}
	return nil
}

func (u *translationUsecase) validateVerseId(ctx context.Context, verseId uint) error {
	verse, err := u.verseRepository.GetById(ctx, verseId)
	if err != nil {
		u.log.Error("failed to get verse for translation", "error", err, "verse_id", verseId)
		return domain.NewInternalError("failed to get verse for translation", err)
	}
	if verse == nil {
		return ErrTranslationNotFound
	}
	return nil
}

// Get a translation by id
func (u *translationUsecase) GetById(ctx context.Context, id uint) (*entity.Translation, error) {
	translation, err := u.translationRepository.GetById(ctx, id)
	if err != nil {
		u.log.Error("failed to get translation by id", "error", err, "translation_id", id)
		return nil, domain.NewInternalError("failed to get translation by id", err)
	}
	if translation == nil {
		return nil, ErrTranslationNotFound
	}
	return translation, nil
}

// List translations for a verse
func (u *translationUsecase) List(ctx context.Context, params portuc.ListParams) (*portuc.PaginatedResult[entity.Translation], error) {
	// apply param default
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}

	offset := (params.Page - 1) * params.Limit

	translations, total, err := u.translationRepository.List(ctx, offset, params.Limit, params.Search)
	if err != nil {
		u.log.Error("failed to list translations", "error", err)
		return nil, domain.NewInternalError("failed to list translations", err)
	}

	totalPages := int(total / uint(params.Limit))
	if uint(total%uint(params.Limit)) > 0 {
		totalPages++
	}
	return &portuc.PaginatedResult[entity.Translation]{
		Data:       translations,
		Total:      int64(total),
		Page:       int(params.Page),
		Limit:      int(params.Limit),
		TotalPages: totalPages,
	}, nil
}

func (u *translationUsecase) GetByVerseId(ctx context.Context, verseId uint) ([]entity.Translation, error) {
	// validate verse id
	if verseId == 0 {
		return nil, ErrTranslationNotFound
	}
	if err := u.validateVerseId(ctx, verseId); err != nil {
		return nil, err
	}

	translations, err := u.translationRepository.GetByVerseId(ctx, verseId)
	if err != nil {
		u.log.Error("failed to get translations by verse id", "error", err, "verse_id", verseId)
		return nil, domain.NewInternalError("failed to get translations by verse id", err)
	}
	return translations, nil
}

// Update a translation
func (u *translationUsecase) Update(ctx context.Context, id uint, input portuc.UpdateTranslationInput) (*entity.Translation, error) {
	translation, err := u.translationRepository.GetById(ctx, id)
	if err != nil {
		u.log.Error("failed to get translation for update", "error", err, "translation_id", id)
		return nil, ErrTranslationNotFound
	}

	// update fields if provided
	if input.VerseID != nil {
		if err := u.validateVerseId(ctx, *input.VerseID); err != nil {
			return nil, err
		}
		translation.VerseID = *input.VerseID
	}
	if input.LanguageCode != nil {
		translation.LanguageCode = *input.LanguageCode
	}
	if input.TranslationText != nil {
		translation.TranslationText = *input.TranslationText
	}
	if input.TranslatorName != nil {
		translation.TranslatorName = input.TranslatorName
	}

	// update translation
	if err := u.translationRepository.Update(ctx, translation); err != nil {
		u.log.Error("failed to update translation", "error", err, "translation_id", id)
		return nil, domain.NewInternalError("failed to update translation", err)
	}
	return translation, nil
}

// GetDropdownData returns dropdown filter data for translations
func (u *translationUsecase) GetDropdownData(ctx context.Context) (*portuc.TranslationDropdownData, error) {
	verses, translatorNames, languageCodes, err := u.translationRepository.GetDropdownData(ctx)
	if err != nil {
		u.log.Error("failed to get translation dropdown data", "error", err)
		return nil, domain.NewInternalError("failed to get translation dropdown data", err)
	}
	return &portuc.TranslationDropdownData{
		Verses:          verses,
		TranslatorNames: translatorNames,
		LanguageCodes:   languageCodes,
	}, nil
}

func (u *translationUsecase) BulkDelete(ctx context.Context, ids []uint) error {
	return u.translationRepository.BulkDelete(ctx, ids)
}

// Delete a translation
func (u *translationUsecase) Delete(ctx context.Context, id uint) error {
	_, err := u.translationRepository.GetById(ctx, id)
	if err != nil {
		u.log.Error("failed to get translation for delete", "error", err, "translation_id", id)
		return ErrTranslationNotFound
	}

	if err := u.translationRepository.Delete(ctx, id); err != nil {
		u.log.Error("failed to delete translation", "error", err, "translation_id", id)
		return domain.NewInternalError("failed to delete translation", err)
	}
	return nil
}
