package bookmark_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/internal/core/usecase/bookmark"
)

// MockBookmarkRepository is a manual mock
type MockBookmarkRepository struct {
	CreateBookmarkFunc                func(ctx context.Context, b *entity.Bookmark) error
	GetBookmarkByIDFunc               func(ctx context.Context, id uint) (*entity.Bookmark, error)
	GetBookmarkByUserIDAndVerseIDFunc func(ctx context.Context, userID uint, verseID uint) (*entity.Bookmark, error)
	ListBookmarksByUserIDFunc         func(ctx context.Context, userID uint, offset, limit int, sort string) ([]entity.Bookmark, int64, error)
	UpdateBookmarkFunc                func(ctx context.Context, b *entity.Bookmark) error
	DeleteBookmarkFunc                func(ctx context.Context, id uint) error
}

func (m *MockBookmarkRepository) CreateBookmark(ctx context.Context, b *entity.Bookmark) error {
	if m.CreateBookmarkFunc != nil {
		return m.CreateBookmarkFunc(ctx, b)
	}
	return nil
}

func (m *MockBookmarkRepository) GetBookmarkByID(ctx context.Context, id uint) (*entity.Bookmark, error) {
	if m.GetBookmarkByIDFunc != nil {
		return m.GetBookmarkByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBookmarkRepository) GetBookmarkByUserIDAndVerseID(ctx context.Context, userID uint, verseID uint) (*entity.Bookmark, error) {
	if m.GetBookmarkByUserIDAndVerseIDFunc != nil {
		return m.GetBookmarkByUserIDAndVerseIDFunc(ctx, userID, verseID)
	}
	return nil, nil
}

func (m *MockBookmarkRepository) ListBookmarksByUserID(ctx context.Context, userID uint, offset, limit int, sort string) ([]entity.Bookmark, int64, error) {
	if m.ListBookmarksByUserIDFunc != nil {
		return m.ListBookmarksByUserIDFunc(ctx, userID, offset, limit, sort)
	}
	return nil, 0, nil
}

func (m *MockBookmarkRepository) UpdateBookmark(ctx context.Context, b *entity.Bookmark) error {
	if m.UpdateBookmarkFunc != nil {
		return m.UpdateBookmarkFunc(ctx, b)
	}
	return nil
}

func (m *MockBookmarkRepository) DeleteBookmark(ctx context.Context, id uint) error {
	if m.DeleteBookmarkFunc != nil {
		return m.DeleteBookmarkFunc(ctx, id)
	}
	return nil
}

// MockVerseRepository is a manual mock
type MockVerseRepository struct {
	GetByIdFunc func(ctx context.Context, id uint) (*entity.Verse, error)
}

func (m *MockVerseRepository) Create(ctx context.Context, v *entity.Verse) error { return nil }
func (m *MockVerseRepository) List(ctx context.Context, filter repository.VerseFilter) ([]entity.Verse, uint, error) {
	return nil, 0, nil
}
func (m *MockVerseRepository) Update(ctx context.Context, v *entity.Verse) error { return nil }
func (m *MockVerseRepository) Delete(ctx context.Context, id uint) error         { return nil }
func (m *MockVerseRepository) BulkDelete(ctx context.Context, ids []uint) error  { return nil }
func (m *MockVerseRepository) GetById(ctx context.Context, id uint) (*entity.Verse, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return nil, nil
}

// MockLogger is a manual mock
type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...any)  {}
func (m *MockLogger) Error(msg string, fields ...any) {}
func (m *MockLogger) Debug(msg string, fields ...any) {}

// Helper string pointer
func strPtr(s string) *string {
	return &s
}

// =============================================================================
// TEST: Create Bookmark
// =============================================================================

func TestBookmarkUseCase_Create_Success(t *testing.T) {
	mockBookmarkRepo := &MockBookmarkRepository{
		GetBookmarkByUserIDAndVerseIDFunc: func(ctx context.Context, userID uint, verseID uint) (*entity.Bookmark, error) {
			return nil, nil // Not found, which is good for creation
		},
		CreateBookmarkFunc: func(ctx context.Context, b *entity.Bookmark) error {
			b.ID = 1
			b.CreatedAt = time.Now()
			return nil
		},
	}

	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1, ArabicText: "Test"}, nil
		},
	}

	mockLogger := &MockLogger{}
	uc := bookmark.NewBookmarkUsecase(mockBookmarkRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateBookmarkInput{
		UserID:  1,
		VerseID: 1,
		Note:    strPtr("Test note"),
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
	if *result.Note != "Test note" {
		t.Errorf("expected Note 'Test note', got '%s'", *result.Note)
	}
}

func TestBookmarkUseCase_Create_VerseNotFound(t *testing.T) {
	mockBookmarkRepo := &MockBookmarkRepository{}
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return nil, errors.New("verse not found")
		},
	}
	mockLogger := &MockLogger{}
	uc := bookmark.NewBookmarkUsecase(mockBookmarkRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateBookmarkInput{
		UserID:  1,
		VerseID: 99,
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestBookmarkUseCase_Create_AlreadyExists(t *testing.T) {
	mockBookmarkRepo := &MockBookmarkRepository{
		GetBookmarkByUserIDAndVerseIDFunc: func(ctx context.Context, userID uint, verseID uint) (*entity.Bookmark, error) {
			return &entity.Bookmark{ID: 1}, nil // Already exists
		},
	}
	mockVerseRepo := &MockVerseRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*entity.Verse, error) {
			return &entity.Verse{ID: 1}, nil
		},
	}
	mockLogger := &MockLogger{}
	uc := bookmark.NewBookmarkUsecase(mockBookmarkRepo, mockVerseRepo, mockLogger)

	input := portuc.CreateBookmarkInput{
		UserID:  1,
		VerseID: 1,
	}

	_, err := uc.Create(context.Background(), input)

	if !errors.Is(err, bookmark.ErrBookmarkAlreadyExists) {
		t.Errorf("expected ErrBookmarkAlreadyExists, got %v", err)
	}
}

// =============================================================================
// TEST: List By User ID
// =============================================================================

func TestBookmarkUseCase_ListByUserID_Success(t *testing.T) {
	mockBookmarkRepo := &MockBookmarkRepository{
		ListBookmarksByUserIDFunc: func(ctx context.Context, userID uint, offset, limit int, sort string) ([]entity.Bookmark, int64, error) {
			return []entity.Bookmark{
				{ID: 1, UserID: 1, VerseID: 1},
				{ID: 2, UserID: 1, VerseID: 2},
			}, 2, nil
		},
	}
	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}
	uc := bookmark.NewBookmarkUsecase(mockBookmarkRepo, mockVerseRepo, mockLogger)

	input := portuc.ListBookmarkInput{
		UserID: 1,
		Page:   1,
		Limit:  10,
	}

	result, err := uc.ListByUserID(context.Background(), input)

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
		t.Errorf("expected 2 items, got %d", len(result.Data))
	}
}

// =============================================================================
// TEST: Update Bookmark
// =============================================================================

func TestBookmarkUseCase_Update_Success(t *testing.T) {
	existingBookmark := &entity.Bookmark{
		ID:      1,
		UserID:  1,
		VerseID: 1,
		Note:    strPtr("Old note"),
	}

	mockBookmarkRepo := &MockBookmarkRepository{
		GetBookmarkByIDFunc: func(ctx context.Context, id uint) (*entity.Bookmark, error) {
			return existingBookmark, nil
		},
		UpdateBookmarkFunc: func(ctx context.Context, b *entity.Bookmark) error {
			return nil
		},
	}
	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}
	uc := bookmark.NewBookmarkUsecase(mockBookmarkRepo, mockVerseRepo, mockLogger)

	input := portuc.UpdateBookmarkInput{
		UserID: 1,
		Note:   strPtr("New note"),
	}

	result, err := uc.Update(context.Background(), 1, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if *result.Note != "New note" {
		t.Errorf("expected 'New note', got '%s'", *result.Note)
	}
}

func TestBookmarkUseCase_Update_Forbidden(t *testing.T) {
	existingBookmark := &entity.Bookmark{
		ID:      1,
		UserID:  2, // Different user
		VerseID: 1,
	}

	mockBookmarkRepo := &MockBookmarkRepository{
		GetBookmarkByIDFunc: func(ctx context.Context, id uint) (*entity.Bookmark, error) {
			return existingBookmark, nil
		},
	}
	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}
	uc := bookmark.NewBookmarkUsecase(mockBookmarkRepo, mockVerseRepo, mockLogger)

	input := portuc.UpdateBookmarkInput{
		UserID: 1, // User 1 trying to edit User 2's bookmark
		Note:   strPtr("Hacked note"),
	}

	_, err := uc.Update(context.Background(), 1, input)

	if !errors.Is(err, bookmark.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

// =============================================================================
// TEST: Delete Bookmark
// =============================================================================

func TestBookmarkUseCase_Delete_Success(t *testing.T) {
	existingBookmark := &entity.Bookmark{
		ID:      1,
		UserID:  1,
		VerseID: 1,
	}

	mockBookmarkRepo := &MockBookmarkRepository{
		GetBookmarkByIDFunc: func(ctx context.Context, id uint) (*entity.Bookmark, error) {
			return existingBookmark, nil
		},
		DeleteBookmarkFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	mockVerseRepo := &MockVerseRepository{}
	mockLogger := &MockLogger{}
	uc := bookmark.NewBookmarkUsecase(mockBookmarkRepo, mockVerseRepo, mockLogger)

	err := uc.Delete(context.Background(), 1, 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
