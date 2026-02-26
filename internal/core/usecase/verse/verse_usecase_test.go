package verse_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/internal/core/usecase/verse"
)

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

// MockChapterRepository is a manual mock for testing
type MockChapterRepository struct {
	CreateChapterFunc       func(ctx context.Context, chapter *entity.Chapter) error
	ListChaptersFunc        func(ctx context.Context, offset, limit int, search string, bookID *uint, title string, category string) ([]entity.Chapter, int64, error)
	GetChaptersByBookIDFunc func(ctx context.Context, bookID uint) ([]entity.Chapter, int64, error)
	GetChapterByIDFunc      func(ctx context.Context, id uint) (*entity.Chapter, error)
	UpdateChapterFunc       func(ctx context.Context, chapter *entity.Chapter) error
	DeleteChapterFunc       func(ctx context.Context, id uint) error
	DeleteChaptersFunc      func(ctx context.Context, ids []uint) error
}

func (m *MockChapterRepository) CreateChapter(ctx context.Context, chapter *entity.Chapter) error {
	if m.CreateChapterFunc != nil {
		return m.CreateChapterFunc(ctx, chapter)
	}
	return nil
}

func (m *MockChapterRepository) ListChapters(ctx context.Context, offset, limit int, search string, bookID *uint, title string, category string) ([]entity.Chapter, int64, error) {
	if m.ListChaptersFunc != nil {
		return m.ListChaptersFunc(ctx, offset, limit, search, bookID, title, category)
	}
	return nil, 0, nil
}

func (m *MockChapterRepository) GetChaptersByBookID(ctx context.Context, bookID uint) ([]entity.Chapter, int64, error) {
	if m.GetChaptersByBookIDFunc != nil {
		return m.GetChaptersByBookIDFunc(ctx, bookID)
	}
	return nil, 0, nil
}

func (m *MockChapterRepository) GetChapterByID(ctx context.Context, id uint) (*entity.Chapter, error) {
	if m.GetChapterByIDFunc != nil {
		return m.GetChapterByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockChapterRepository) UpdateChapter(ctx context.Context, chapter *entity.Chapter) error {
	if m.UpdateChapterFunc != nil {
		return m.UpdateChapterFunc(ctx, chapter)
	}
	return nil
}

func (m *MockChapterRepository) DeleteChapter(ctx context.Context, id uint) error {
	if m.DeleteChapterFunc != nil {
		return m.DeleteChapterFunc(ctx, id)
	}
	return nil
}

func (m *MockChapterRepository) DeleteChapters(ctx context.Context, ids []uint) error {
	if m.DeleteChaptersFunc != nil {
		return m.DeleteChaptersFunc(ctx, ids)
	}
	return nil
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
// TEST: Create Verse
// =============================================================================

func TestVerseUseCase_Create_Success(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		CreateFunc: func(ctx context.Context, v *entity.Verse) error {
			v.ID = 1
			v.CreatedAt = time.Now()
			v.UpdatedAt = time.Now()
			return nil
		},
	}

	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1, Title: "Test Chapter"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.CreateVerseInput{
		ChapterID:       1,
		VerseNumber:     1,
		ArabicText:      "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ",
		Transliteration: strPtr("Bismillahir Rahmanir Rahim"),
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
	if result.ChapterID != 1 {
		t.Errorf("expected ChapterID 1, got %d", result.ChapterID)
	}
	if result.ArabicText != input.ArabicText {
		t.Errorf("expected ArabicText '%s', got '%s'", input.ArabicText, result.ArabicText)
	}
}

func TestVerseUseCase_Create_ChapterNotFound(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{}

	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return nil, nil // Chapter not found
		},
	}

	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.CreateVerseInput{
		ChapterID:   999,
		VerseNumber: 1,
		ArabicText:  "Test Arabic Text",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestVerseUseCase_Create_InvalidChapterID(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{}
	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.CreateVerseInput{
		ChapterID:   0,
		VerseNumber: 1,
		ArabicText:  "Test Arabic Text",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid chapter ID, got nil")
	}
}

func TestVerseUseCase_Create_InvalidVerseNumber(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{}

	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1, Title: "Test Chapter"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.CreateVerseInput{
		ChapterID:   1,
		VerseNumber: 0,
		ArabicText:  "Test Arabic Text",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid verse number, got nil")
	}
}

func TestVerseUseCase_Create_EmptyArabicText(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{}

	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1, Title: "Test Chapter"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.CreateVerseInput{
		ChapterID:   1,
		VerseNumber: 1,
		ArabicText:  "",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty arabic text, got nil")
	}
}

func TestVerseUseCase_Create_RepositoryError(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		CreateFunc: func(ctx context.Context, v *entity.Verse) error {
			return errors.New("db error")
		},
	}

	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1, Title: "Test Chapter"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.CreateVerseInput{
		ChapterID:   1,
		VerseNumber: 1,
		ArabicText:  "Test Arabic Text",
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// =============================================================================
// TEST: List Verses
// =============================================================================

func TestVerseUseCase_List_Success(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		ListFunc: func(ctx context.Context, filter repository.VerseFilter) ([]entity.Verse, uint, error) {
			return []entity.Verse{
				{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Test 1"},
				{ID: 2, ChapterID: 1, VerseNumber: 2, ArabicText: "Test 2"},
			}, 2, nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	params := portuc.ListParams{
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
		t.Errorf("expected 2 verses, got %d", len(result.Data))
	}
}

func TestVerseUseCase_List_DefaultParams(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		ListFunc: func(ctx context.Context, filter repository.VerseFilter) ([]entity.Verse, uint, error) {
			// Verify default values were applied
			if filter.Offset != 0 {
				t.Errorf("expected offset 0, got %d", filter.Offset)
			}
			if filter.Limit != 20 {
				t.Errorf("expected limit 20, got %d", filter.Limit)
			}
			return []entity.Verse{}, 0, nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	params := portuc.ListParams{
		Page:  0, // Should default to 1
		Limit: 0, // Should default to 20
	}

	_, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestVerseUseCase_List_WithSearch(t *testing.T) {
	searchTerm := "bismillah"
	mockVerseRepo := &MockVerseRepository{
		ListFunc: func(ctx context.Context, filter repository.VerseFilter) ([]entity.Verse, uint, error) {
			if filter.Search != searchTerm {
				t.Errorf("expected search '%s', got '%s'", searchTerm, filter.Search)
			}
			return []entity.Verse{
				{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Bismillah"},
			}, 1, nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	params := portuc.ListParams{
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

func TestVerseUseCase_List_RepositoryError(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		ListFunc: func(ctx context.Context, filter repository.VerseFilter) ([]entity.Verse, uint, error) {
			return nil, 0, errors.New("db error")
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	params := portuc.ListParams{
		Page:  1,
		Limit: 10,
	}

	_, err := uc.List(context.Background(), params)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestVerseUseCase_List_Pagination(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		ListFunc: func(ctx context.Context, filter repository.VerseFilter) ([]entity.Verse, uint, error) {
			return []entity.Verse{
				{ID: 1, ChapterID: 1, VerseNumber: 1, ArabicText: "Test"},
			}, 25, nil // Total 25 items
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	params := portuc.ListParams{
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
// TEST: GetById
// =============================================================================

func TestVerseUseCase_GetById_Success(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			if id == 1 {
				return &entity.Verse{
					ID:          1,
					ChapterID:   1,
					VerseNumber: 1,
					ArabicText:  "Test Arabic Text",
				}, nil
			}
			return nil, nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

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

func TestVerseUseCase_GetById_NotFound(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	_, err := uc.GetById(context.Background(), 999)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestVerseUseCase_GetById_RepositoryError(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, errors.New("db error")
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	_, err := uc.GetById(context.Background(), 1)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// =============================================================================
// TEST: Update Verse
// =============================================================================

func TestVerseUseCase_Update_Success(t *testing.T) {
	existingVerse := &entity.Verse{
		ID:          1,
		ChapterID:   1,
		VerseNumber: 1,
		ArabicText:  "Old Arabic Text",
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			if id == 1 {
				return existingVerse, nil
			}
			return nil, nil
		},
		UpdateFunc: func(ctx context.Context, v *entity.Verse) error {
			return nil
		},
	}

	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 2, Title: "New Chapter"}, nil
		},
	}

	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.UpdateVerseInput{
		ChapterID:       uintPtr(2),
		VerseNumber:     uintPtr(5),
		ArabicText:      strPtr("New Arabic Text"),
		Transliteration: strPtr("New Transliteration"),
	}

	result, err := uc.Update(context.Background(), 1, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.ChapterID != 2 {
		t.Errorf("expected ChapterID 2, got %d", result.ChapterID)
	}
	if result.VerseNumber != 5 {
		t.Errorf("expected VerseNumber 5, got %d", result.VerseNumber)
	}
	if result.ArabicText != "New Arabic Text" {
		t.Errorf("expected ArabicText 'New Arabic Text', got '%s'", result.ArabicText)
	}
}

func TestVerseUseCase_Update_NotFound(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, errors.New("not found")
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.UpdateVerseInput{
		ArabicText: strPtr("New Text"),
	}

	_, err := uc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestVerseUseCase_Update_InvalidChapter(t *testing.T) {
	existingVerse := &entity.Verse{
		ID:          1,
		ChapterID:   1,
		VerseNumber: 1,
		ArabicText:  "Old Arabic Text",
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return existingVerse, nil
		},
	}

	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return nil, nil // Chapter not found
		},
	}

	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.UpdateVerseInput{
		ChapterID: uintPtr(999),
	}

	_, err := uc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for invalid chapter, got nil")
	}
}

func TestVerseUseCase_Update_InvalidVerseNumber(t *testing.T) {
	existingVerse := &entity.Verse{
		ID:          1,
		ChapterID:   1,
		VerseNumber: 1,
		ArabicText:  "Old Arabic Text",
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return existingVerse, nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.UpdateVerseInput{
		VerseNumber: uintPtr(0),
	}

	_, err := uc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for invalid verse number, got nil")
	}
}

func TestVerseUseCase_Update_EmptyArabicText(t *testing.T) {
	existingVerse := &entity.Verse{
		ID:          1,
		ChapterID:   1,
		VerseNumber: 1,
		ArabicText:  "Old Arabic Text",
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return existingVerse, nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.UpdateVerseInput{
		ArabicText: strPtr(""),
	}

	_, err := uc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for empty arabic text, got nil")
	}
}

func TestVerseUseCase_Update_RepositoryError(t *testing.T) {
	existingVerse := &entity.Verse{
		ID:          1,
		ChapterID:   1,
		VerseNumber: 1,
		ArabicText:  "Old Arabic Text",
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return existingVerse, nil
		},
		UpdateFunc: func(ctx context.Context, v *entity.Verse) error {
			return errors.New("db error")
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	input := portuc.UpdateVerseInput{
		Transliteration: strPtr("New Transliteration"),
	}

	_, err := uc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestVerseUseCase_Update_PartialUpdate(t *testing.T) {
	existingVerse := &entity.Verse{
		ID:              1,
		ChapterID:       1,
		VerseNumber:     1,
		ArabicText:      "Old Arabic Text",
		Transliteration: strPtr("Old Transliteration"),
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return existingVerse, nil
		},
		UpdateFunc: func(ctx context.Context, v *entity.Verse) error {
			return nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	// Only update transliteration
	input := portuc.UpdateVerseInput{
		Transliteration: strPtr("New Transliteration Only"),
	}

	result, err := uc.Update(context.Background(), 1, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	// Original values should remain unchanged
	if result.ChapterID != 1 {
		t.Errorf("expected ChapterID unchanged at 1, got %d", result.ChapterID)
	}
	if result.VerseNumber != 1 {
		t.Errorf("expected VerseNumber unchanged at 1, got %d", result.VerseNumber)
	}
	if result.ArabicText != "Old Arabic Text" {
		t.Errorf("expected ArabicText unchanged, got '%s'", result.ArabicText)
	}
	// Transliteration should be updated
	if *result.Transliteration != "New Transliteration Only" {
		t.Errorf("expected Transliteration 'New Transliteration Only', got '%s'", *result.Transliteration)
	}
}

// =============================================================================
// TEST: Delete Verse
// =============================================================================

func TestVerseUseCase_Delete_Success(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			if id == 1 {
				return &entity.Verse{ID: 1}, nil
			}
			return nil, nil
		},
		DeleteFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	err := uc.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestVerseUseCase_Delete_NotFound(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, errors.New("not found")
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	err := uc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestVerseUseCase_Delete_RepositoryError(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1}, nil
		},
		DeleteFunc: func(ctx context.Context, id uint) error {
			return errors.New("db error")
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	err := uc.Delete(context.Background(), 1)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// =============================================================================
// TEST: BulkDelete Verses
// =============================================================================

func TestVerseUseCase_BulkDelete_Success(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		BulkDeleteFunc: func(ctx context.Context, ids []uint) error {
			return nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	err := uc.BulkDelete(context.Background(), []uint{1, 2, 3})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestVerseUseCase_BulkDelete_EmptyIDs(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		BulkDeleteFunc: func(ctx context.Context, ids []uint) error {
			t.Error("BulkDelete should not be called")
			return nil
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	err := uc.BulkDelete(context.Background(), []uint{})

	if err != nil {
		t.Errorf("expected no error for empty IDs, got %v", err)
	}
}

func TestVerseUseCase_BulkDelete_RepositoryError(t *testing.T) {
	mockVerseRepo := &MockVerseRepository{
		BulkDeleteFunc: func(ctx context.Context, ids []uint) error {
			return errors.New("db error")
		},
	}

	mockChapterRepo := &MockChapterRepository{}
	mockLogger := &MockLogger{}

	uc := verse.NewVerseUsecase(mockVerseRepo, mockChapterRepo, mockLogger)

	err := uc.BulkDelete(context.Background(), []uint{1, 2})

	if err == nil {
		t.Error("expected error, got nil")
	}
}
