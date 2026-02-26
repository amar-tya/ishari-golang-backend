package translation_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/internal/core/usecase/translation"
)

// MockTranslationRepository is a manual mock for testing
type MockTranslationRepository struct {
	CreateFunc          func(ctx context.Context, translation *entity.Translation) error
	ListFunc            func(ctx context.Context, filter repository.TranslationListFilter) ([]entity.Translation, uint, error)
	UpdateFunc          func(ctx context.Context, translation *entity.Translation) error
	DeleteFunc          func(ctx context.Context, id uint) error
	GetByIdFunc         func(ctx context.Context, id uint) (*entity.Translation, error)
	GetByVerseIdFunc    func(ctx context.Context, verseId uint) ([]entity.Translation, error)
	GetDropdownDataFunc func(ctx context.Context) ([]entity.Verse, []string, []string, error)
	BulkDeleteFunc      func(ctx context.Context, ids []uint) error
}

func (m *MockTranslationRepository) Create(ctx context.Context, t *entity.Translation) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, t)
	}
	return nil
}

func (m *MockTranslationRepository) List(ctx context.Context, filter repository.TranslationListFilter) ([]entity.Translation, uint, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filter)
	}
	return nil, 0, nil
}

func (m *MockTranslationRepository) Update(ctx context.Context, t *entity.Translation) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, t)
	}
	return nil
}

func (m *MockTranslationRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTranslationRepository) GetById(ctx context.Context, id uint) (*entity.Translation, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTranslationRepository) GetByVerseId(ctx context.Context, verseId uint) ([]entity.Translation, error) {
	if m.GetByVerseIdFunc != nil {
		return m.GetByVerseIdFunc(ctx, verseId)
	}
	return nil, nil
}

func (m *MockTranslationRepository) GetDropdownData(ctx context.Context) ([]entity.Verse, []string, []string, error) {
	if m.GetDropdownDataFunc != nil {
		return m.GetDropdownDataFunc(ctx)
	}
	return nil, nil, nil, nil
}

func (m *MockTranslationRepository) BulkDelete(ctx context.Context, ids []uint) error {
	if m.BulkDeleteFunc != nil {
		return m.BulkDeleteFunc(ctx, ids)
	}
	return nil
}

// MockVerseRepository is a manual mock for testing
type MockVerseRepository struct {
	CreateFunc     func(ctx context.Context, verse *entity.Verse) error
	ListFunc       func(ctx context.Context, filter repository.VerseFilter) ([]entity.Verse, uint, error)
	UpdateFunc     func(ctx context.Context, verse *entity.Verse) error
	DeleteFunc     func(ctx context.Context, id uint) error
	BulkDeleteFunc func(ctx context.Context, ids []uint) error
	GetByIdFunc    func(ctx context.Context, id uint) (*entity.Verse, error)
}

func (m *MockVerseRepository) Create(ctx context.Context, v *entity.Verse) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, v)
	}
	return nil
}

func (m *MockVerseRepository) List(ctx context.Context, filter repository.VerseFilter) ([]entity.Verse, uint, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filter)
	}
	return nil, 0, nil
}

func (m *MockVerseRepository) Update(ctx context.Context, v *entity.Verse) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, v)
	}
	return nil
}

func (m *MockVerseRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockVerseRepository) BulkDelete(ctx context.Context, ids []uint) error {
	if m.BulkDeleteFunc != nil {
		return m.BulkDeleteFunc(ctx, ids)
	}
	return nil
}

func (m *MockVerseRepository) GetById(ctx context.Context, id uint) (*entity.Verse, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return nil, nil
}

// MockLogger is a manual mock for testing
type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...any)  {}
func (m *MockLogger) Error(msg string, fields ...any) {}

// Helper function to create a pointer to a string
func strPtr(s string) *string {
	return &s
}

// Helper function to create a pointer to uint
func uintPtr(u uint) *uint {
	return &u
}

// =============================================================================
// TEST: Create Translation
// =============================================================================

func TestTranslationUseCase_Create_Success(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		CreateFunc: func(ctx context.Context, tr *entity.Translation) error {
			tr.ID = 1
			tr.CreatedAt = time.Now()
			tr.UpdatedAt = time.Now()
			return nil
		},
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			if id == 1 {
				return &entity.Translation{
					ID:              1,
					VerseID:         1,
					LanguageCode:    "en",
					TranslationText: "In the name of God, the Most Gracious, the Most Merciful",
					TranslatorName:  strPtr("Test Translator"),
				}, nil
			}
			return nil, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Test Arabic"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateTranslationInput{
		VerseID:         1,
		LanguageCode:    "en",
		TranslationText: "In the name of God, the Most Gracious, the Most Merciful",
		TranslatorName:  strPtr("Test Translator"),
	}

	result, err := uc.Create(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.VerseID != 1 {
		t.Errorf("expected VerseID 1, got %d", result.VerseID)
	}
	if result.LanguageCode != input.LanguageCode {
		t.Errorf("expected LanguageCode '%s', got '%s'", input.LanguageCode, result.LanguageCode)
	}
	if result.TranslationText != input.TranslationText {
		t.Errorf("expected TranslationText '%s', got '%s'", input.TranslationText, result.TranslationText)
	}
}

func TestTranslationUseCase_Create_VerseNotFound(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, nil // Verse not found
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateTranslationInput{
		VerseID:         999,
		LanguageCode:    "en",
		TranslationText: "Test Translation",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_Create_InvalidVerseID(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{}
	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateTranslationInput{
		VerseID:         0,
		LanguageCode:    "en",
		TranslationText: "Test Translation",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid verse ID, got nil")
	}
}

func TestTranslationUseCase_Create_EmptyLanguageCode(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Test Arabic"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateTranslationInput{
		VerseID:         1,
		LanguageCode:    "",
		TranslationText: "Test Translation",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty language code, got nil")
	}
}

func TestTranslationUseCase_Create_EmptyTranslationText(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Test Arabic"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateTranslationInput{
		VerseID:         1,
		LanguageCode:    "en",
		TranslationText: "",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty translation text, got nil")
	}
}

func TestTranslationUseCase_Create_RepositoryError(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		CreateFunc: func(ctx context.Context, tr *entity.Translation) error {
			return errors.New("db error")
		},
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Test Arabic"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateTranslationInput{
		VerseID:         1,
		LanguageCode:    "en",
		TranslationText: "Test Translation",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_Create_VerseRepositoryError(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, errors.New("db error")
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateTranslationInput{
		VerseID:         1,
		LanguageCode:    "en",
		TranslationText: "Test Translation",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// =============================================================================
// TEST: GetById Translation
// =============================================================================

func TestTranslationUseCase_GetById_Success(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			if id == 1 {
				return &entity.Translation{
					ID:              1,
					VerseID:         1,
					LanguageCode:    "en",
					TranslationText: "Test Translation",
					TranslatorName:  strPtr("Test Translator"),
				}, nil
			}
			return nil, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	result, err := uc.GetById(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestTranslationUseCase_GetById_NotFound(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return nil, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	_, err := uc.GetById(context.Background(), 999)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_GetById_RepositoryError(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return nil, errors.New("db error")
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	_, err := uc.GetById(context.Background(), 1)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// =============================================================================
// TEST: List Translations
// =============================================================================

func TestTranslationUseCase_List_Success(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		ListFunc: func(ctx context.Context, filter repository.TranslationListFilter) ([]entity.Translation, uint, error) {
			return []entity.Translation{
				{ID: 1, VerseID: 1, LanguageCode: "en", TranslationText: "Test 1"},
				{ID: 2, VerseID: 1, LanguageCode: "id", TranslationText: "Test 2"},
			}, 2, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	params := portuc.TranslationListParams{
		Page:  1,
		Limit: 10,
	}

	result, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
	if len(result.Data) != 2 {
		t.Errorf("expected 2 translations, got %d", len(result.Data))
	}
}

func TestTranslationUseCase_List_DefaultParams(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		ListFunc: func(ctx context.Context, filter repository.TranslationListFilter) ([]entity.Translation, uint, error) {
			// Verify default values were applied
			if filter.Offset != 0 {
				t.Errorf("expected offset 0, got %d", filter.Offset)
			}
			if filter.Limit != 20 {
				t.Errorf("expected limit 20, got %d", filter.Limit)
			}
			return []entity.Translation{}, 0, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	params := portuc.TranslationListParams{
		Page:  0, // Should default to 1
		Limit: 0, // Should default to 20
	}

	_, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTranslationUseCase_List_WithSearch(t *testing.T) {
	searchTerm := "bismillah"
	mockTransRepo := &MockTranslationRepository{
		ListFunc: func(ctx context.Context, filter repository.TranslationListFilter) ([]entity.Translation, uint, error) {
			if filter.Search != searchTerm {
				t.Errorf("expected search '%s', got '%s'", searchTerm, filter.Search)
			}
			return []entity.Translation{
				{ID: 1, VerseID: 1, LanguageCode: "en", TranslationText: "In the name of God"},
			}, 1, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	params := portuc.TranslationListParams{
		Page:   1,
		Limit:  10,
		Search: searchTerm,
	}

	result, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
}

func TestTranslationUseCase_List_RepositoryError(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		ListFunc: func(ctx context.Context, filter repository.TranslationListFilter) ([]entity.Translation, uint, error) {
			return nil, 0, errors.New("db error")
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	params := portuc.TranslationListParams{
		Page:  1,
		Limit: 10,
	}

	_, err := uc.List(context.Background(), params)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_List_Pagination(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		ListFunc: func(ctx context.Context, filter repository.TranslationListFilter) ([]entity.Translation, uint, error) {
			return []entity.Translation{
				{ID: 1, VerseID: 1, LanguageCode: "en", TranslationText: "Test"},
			}, 25, nil // Total 25 items
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	params := portuc.TranslationListParams{
		Page:  1,
		Limit: 10,
	}

	result, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.TotalPages != 3 { // 25/10 = 2.5, rounded up = 3
		t.Errorf("expected total pages 3, got %d", result.TotalPages)
	}
}

// =============================================================================
// TEST: GetByVerseId Translations
// =============================================================================

func TestTranslationUseCase_GetByVerseId_Success(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByVerseIdFunc: func(ctx context.Context, verseId uint) ([]entity.Translation, error) {
			return []entity.Translation{
				{ID: 1, VerseID: verseId, LanguageCode: "en", TranslationText: "English Translation"},
				{ID: 2, VerseID: verseId, LanguageCode: "id", TranslationText: "Terjemahan Indonesia"},
			}, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Test Arabic"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	result, err := uc.GetByVerseId(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 translations, got %d", len(result))
	}
}

func TestTranslationUseCase_GetByVerseId_InvalidVerseID(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{}
	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	_, err := uc.GetByVerseId(context.Background(), 0)

	if err == nil {
		t.Error("expected error for invalid verse ID, got nil")
	}
}

func TestTranslationUseCase_GetByVerseId_VerseNotFound(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, nil // Verse not found
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	_, err := uc.GetByVerseId(context.Background(), 999)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_GetByVerseId_RepositoryError(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByVerseIdFunc: func(ctx context.Context, verseId uint) ([]entity.Translation, error) {
			return nil, errors.New("db error")
		},
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Test Arabic"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	_, err := uc.GetByVerseId(context.Background(), 1)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// =============================================================================
// TEST: Update Translation
// =============================================================================

func TestTranslationUseCase_Update_Success(t *testing.T) {
	existingTranslation := &entity.Translation{
		ID:              1,
		VerseID:         1,
		LanguageCode:    "en",
		TranslationText: "Old Translation",
		TranslatorName:  strPtr("Old Translator"),
	}

	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			if id == 1 {
				return existingTranslation, nil
			}
			return nil, nil
		},
		UpdateFunc: func(ctx context.Context, t *entity.Translation) error {
			return nil
		},
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 2, ChapterID: 1, VerseNumber: 2, ArabicText: "Test Arabic"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.UpdateTranslationInput{
		VerseID:         uintPtr(2),
		LanguageCode:    strPtr("id"),
		TranslationText: strPtr("New Translation"),
		TranslatorName:  strPtr("New Translator"),
	}

	result, err := uc.Update(context.Background(), 1, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.VerseID != 2 {
		t.Errorf("expected VerseID 2, got %d", result.VerseID)
	}
	if result.LanguageCode != "id" {
		t.Errorf("expected LanguageCode 'id', got '%s'", result.LanguageCode)
	}
	if result.TranslationText != "New Translation" {
		t.Errorf("expected TranslationText 'New Translation', got '%s'", result.TranslationText)
	}
}

func TestTranslationUseCase_Update_NotFound(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return nil, errors.New("not found")
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.UpdateTranslationInput{
		TranslationText: strPtr("New Text"),
	}

	_, err := uc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_Update_InvalidVerse(t *testing.T) {
	existingTranslation := &entity.Translation{
		ID:              1,
		VerseID:         1,
		LanguageCode:    "en",
		TranslationText: "Old Translation",
	}

	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return existingTranslation, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, nil // Verse not found
		},
	}

	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.UpdateTranslationInput{
		VerseID: uintPtr(999),
	}

	_, err := uc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for invalid verse, got nil")
	}
}

func TestTranslationUseCase_Update_RepositoryError(t *testing.T) {
	existingTranslation := &entity.Translation{
		ID:              1,
		VerseID:         1,
		LanguageCode:    "en",
		TranslationText: "Old Translation",
	}

	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return existingTranslation, nil
		},
		UpdateFunc: func(ctx context.Context, t *entity.Translation) error {
			return errors.New("db error")
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	input := portuc.UpdateTranslationInput{
		TranslatorName: strPtr("New Translator"),
	}

	_, err := uc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_Update_PartialUpdate(t *testing.T) {
	existingTranslation := &entity.Translation{
		ID:              1,
		VerseID:         1,
		LanguageCode:    "en",
		TranslationText: "Old Translation",
		TranslatorName:  strPtr("Old Translator"),
	}

	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return existingTranslation, nil
		},
		UpdateFunc: func(ctx context.Context, t *entity.Translation) error {
			return nil
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	// Only update translator name
	input := portuc.UpdateTranslationInput{
		TranslatorName: strPtr("New Translator Only"),
	}

	result, err := uc.Update(context.Background(), 1, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	// Original values should remain unchanged
	if result.VerseID != 1 {
		t.Errorf("expected VerseID unchanged at 1, got %d", result.VerseID)
	}
	if result.LanguageCode != "en" {
		t.Errorf("expected LanguageCode unchanged at 'en', got '%s'", result.LanguageCode)
	}
	if result.TranslationText != "Old Translation" {
		t.Errorf("expected TranslationText unchanged, got '%s'", result.TranslationText)
	}
	// TranslatorName should be updated
	if *result.TranslatorName != "New Translator Only" {
		t.Errorf("expected TranslatorName 'New Translator Only', got '%s'", *result.TranslatorName)
	}
}

// =============================================================================
// TEST: Delete Translation
// =============================================================================

func TestTranslationUseCase_Delete_Success(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return &entity.Translation{
				ID:              1,
				VerseID:         1,
				LanguageCode:    "en",
				TranslationText: "Test Translation",
			}, nil
		},
		DeleteFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	err := uc.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTranslationUseCase_Delete_NotFound(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return nil, errors.New("not found")
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	err := uc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_Delete_RepositoryError(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Translation, error) {
			return &entity.Translation{
				ID:              1,
				VerseID:         1,
				LanguageCode:    "en",
				TranslationText: "Test Translation",
			}, nil
		},
		DeleteFunc: func(ctx context.Context, id uint) error {
			return errors.New("db error")
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	err := uc.Delete(context.Background(), 1)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTranslationUseCase_List_WithFilters(t *testing.T) {
	mockTransRepo := &MockTranslationRepository{
		ListFunc: func(ctx context.Context, filter repository.TranslationListFilter) ([]entity.Translation, uint, error) {
			if filter.VerseID != 1 {
				t.Errorf("expected VerseID 1, got %d", filter.VerseID)
			}
			if filter.TranslatorName != "Test Translator" {
				t.Errorf("expected TranslatorName 'Test Translator', got '%s'", filter.TranslatorName)
			}
			if filter.LanguageCode != "en" {
				t.Errorf("expected LanguageCode 'en', got '%s'", filter.LanguageCode)
			}
			return []entity.Translation{
				{ID: 1, VerseID: 1, LanguageCode: "en", TranslationText: "In the name of God", TranslatorName: strPtr("Test Translator")},
			}, 1, nil
		},
	}

	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}

	uc := translation.NewTranslationUsecase(mockTransRepo, mockVerseRepo, mockLogger)

	params := portuc.TranslationListParams{
		Page:           1,
		Limit:          10,
		VerseID:        1,
		TranslatorName: "Test Translator",
		LanguageCode:   "en",
	}

	result, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
}
