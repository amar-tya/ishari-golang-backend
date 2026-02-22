package chapter_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/core/entity"
	portuc "ishari-backend/internal/core/port/usecase"
	"ishari-backend/internal/core/usecase/chapter"
)

// MockChapterRepository is a manual mock for ChapterRepository
type MockChapterRepository struct {
	CreateChapterFunc       func(ctx context.Context, ch *entity.Chapter) error
	ListChaptersFunc        func(ctx context.Context, offset, limit int, search string, bookID *uint, title string, category string) ([]entity.Chapter, int64, error)
	GetChaptersByBookIDFunc func(ctx context.Context, bookID uint) ([]entity.Chapter, int64, error)
	GetChapterByIDFunc      func(ctx context.Context, id uint) (*entity.Chapter, error)
	UpdateChapterFunc       func(ctx context.Context, ch *entity.Chapter) error
	DeleteChapterFunc       func(ctx context.Context, id uint) error
	DeleteChaptersFunc      func(ctx context.Context, ids []uint) error
}

func (m *MockChapterRepository) CreateChapter(ctx context.Context, ch *entity.Chapter) error {
	if m.CreateChapterFunc != nil {
		return m.CreateChapterFunc(ctx, ch)
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

func (m *MockChapterRepository) UpdateChapter(ctx context.Context, ch *entity.Chapter) error {
	if m.UpdateChapterFunc != nil {
		return m.UpdateChapterFunc(ctx, ch)
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

// MockBookRepository is a manual mock for BookRepository
type MockBookRepository struct {
	GetByIdFunc func(ctx context.Context, id int64) (*entity.Book, error)
	ListFunc    func(ctx context.Context, offset, limit int, search string) ([]entity.Book, int64, error)
	CreateFunc  func(ctx context.Context, book *entity.Book) error
	EditFunc    func(ctx context.Context, book *entity.Book) error
	DeleteFunc  func(ctx context.Context, id int64) error
}

func (m *MockBookRepository) GetById(ctx context.Context, id int64) (*entity.Book, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBookRepository) List(ctx context.Context, offset, limit int, search string) ([]entity.Book, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, offset, limit, search)
	}
	return nil, 0, nil
}

func (m *MockBookRepository) Create(ctx context.Context, book *entity.Book) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, book)
	}
	return nil
}

func (m *MockBookRepository) Edit(ctx context.Context, book *entity.Book) error {
	if m.EditFunc != nil {
		return m.EditFunc(ctx, book)
	}
	return nil
}

func (m *MockBookRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockLogger is a manual mock for Logger
type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...any)  {}
func (m *MockLogger) Error(msg string, fields ...any) {}

// Helper function to create description pointer
func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}

// ==================== Create Tests ====================

func TestChapterUsecase_Create_Success(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		CreateChapterFunc: func(ctx context.Context, ch *entity.Chapter) error {
			ch.ID = 1
			ch.CreatedAt = time.Now()
			ch.UpdatedAt = time.Now()
			return nil
		},
	}
	mockBookRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			return &entity.Book{ID: 1, Title: "Test Book"}, nil
		},
	}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	input := portuc.CreateChapterInput{
		BookID:        1,
		ChapterNumber: 1,
		Title:         "Genesis",
		Category:      "Old Testament",
		TotalVerses:   31,
		Description:   stringPtr("First chapter"),
	}

	result, err := uc.Create(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.Title != "Genesis" {
		t.Errorf("expected title 'Genesis', got %s", result.Title)
	}
}

func TestChapterUsecase_Create_BookNotFound(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{}
	mockBookRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			return nil, nil // Book not found
		},
	}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	input := portuc.CreateChapterInput{
		BookID:        999,
		ChapterNumber: 1,
		Title:         "Test",
		Category:      "Test",
		TotalVerses:   1,
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for non-existent book, got nil")
	}
}

func TestChapterUsecase_Create_InvalidChapterNumber(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{}
	mockBookRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			return &entity.Book{ID: 1}, nil
		},
	}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	input := portuc.CreateChapterInput{
		BookID:        1,
		ChapterNumber: 0, // Invalid: must be > 0
		Title:         "Test",
		Category:      "Test",
		TotalVerses:   1,
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid chapter number, got nil")
	}
}

func TestChapterUsecase_Create_EmptyTitle(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{}
	mockBookRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			return &entity.Book{ID: 1}, nil
		},
	}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	input := portuc.CreateChapterInput{
		BookID:        1,
		ChapterNumber: 1,
		Title:         "", // Invalid: empty
		Category:      "Test",
		TotalVerses:   1,
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty title, got nil")
	}
}

func TestChapterUsecase_Create_EmptyCategory(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{}
	mockBookRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			return &entity.Book{ID: 1}, nil
		},
	}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	input := portuc.CreateChapterInput{
		BookID:        1,
		ChapterNumber: 1,
		Title:         "Test",
		Category:      "", // Invalid: empty
		TotalVerses:   1,
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty category, got nil")
	}
}

func TestChapterUsecase_Create_InvalidTotalVerses(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{}
	mockBookRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			return &entity.Book{ID: 1}, nil
		},
	}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	input := portuc.CreateChapterInput{
		BookID:        1,
		ChapterNumber: 1,
		Title:         "Test",
		Category:      "Test",
		TotalVerses:   0, // Invalid: must be > 0
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid total verses, got nil")
	}
}

func TestChapterUsecase_Create_RepositoryError(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		CreateChapterFunc: func(ctx context.Context, ch *entity.Chapter) error {
			return errors.New("database error")
		},
	}
	mockBookRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			return &entity.Book{ID: 1}, nil
		},
	}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	input := portuc.CreateChapterInput{
		BookID:        1,
		ChapterNumber: 1,
		Title:         "Test",
		Category:      "Test",
		TotalVerses:   1,
	}

	_, err := uc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error from repository, got nil")
	}
}

// ==================== List Tests ====================

func TestChapterUsecase_List_Success(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		ListChaptersFunc: func(ctx context.Context, offset, limit int, search string, bookID *uint, title string, category string) ([]entity.Chapter, int64, error) {
			return []entity.Chapter{
				{ID: 1, Title: "Chapter 1"},
				{ID: 2, Title: "Chapter 2"},
			}, 2, nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	params := portuc.ListChapterInput{
		Page:  1,
		Limit: 20,
	}

	result, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
	if len(result.Data) != 2 {
		t.Errorf("expected 2 chapters, got %d", len(result.Data))
	}
}

func TestChapterUsecase_List_WithFilters(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		ListChaptersFunc: func(ctx context.Context, offset, limit int, search string, bookID *uint, title string, category string) ([]entity.Chapter, int64, error) {
			if bookID == nil || *bookID != 1 {
				t.Errorf("expected bookID 1, got %v", bookID)
			}
			if title != "Test" {
				t.Errorf("expected title 'Test', got %s", title)
			}
			if category != "Test Cat" {
				t.Errorf("expected category 'Test Cat', got %s", category)
			}
			return []entity.Chapter{
				{ID: 1, Title: "Chapter 1"},
			}, 1, nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	id := uint(1)
	params := portuc.ListChapterInput{
		Page:     1,
		Limit:    20,
		BookID:   &id,
		Title:    "Test",
		Category: "Test Cat",
	}

	_, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestChapterUsecase_List_DefaultValues(t *testing.T) {
	var capturedOffset, capturedLimit int
	mockChapterRepo := &MockChapterRepository{
		ListChaptersFunc: func(ctx context.Context, offset, limit int, search string, bookID *uint, title string, category string) ([]entity.Chapter, int64, error) {
			capturedOffset = offset
			capturedLimit = limit
			return []entity.Chapter{}, 0, nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	params := portuc.ListChapterInput{
		Page:  0,  // Should default to 1
		Limit: -1, // Should default to 20
	}

	_, err := uc.List(context.Background(), params)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if capturedOffset != 0 {
		t.Errorf("expected offset 0, got %d", capturedOffset)
	}
	if capturedLimit != 20 {
		t.Errorf("expected limit 20, got %d", capturedLimit)
	}
}

func TestChapterUsecase_List_RepositoryError(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		ListChaptersFunc: func(ctx context.Context, offset, limit int, search string, bookID *uint, title string, category string) ([]entity.Chapter, int64, error) {
			return nil, 0, errors.New("database error")
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	params := portuc.ListChapterInput{Page: 1, Limit: 20}

	_, err := uc.List(context.Background(), params)

	if err == nil {
		t.Error("expected error from repository, got nil")
	}
}

// ==================== GetByID Tests ====================

func TestChapterUsecase_GetByID_Success(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1, Title: "Test Chapter"}, nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	result, err := uc.GetByID(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestChapterUsecase_GetByID_NotFound(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return nil, errors.New("not found")
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	_, err := uc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent chapter, got nil")
	}
}

// ==================== GetByBookID Tests ====================

func TestChapterUsecase_GetByBookID_Success(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChaptersByBookIDFunc: func(ctx context.Context, bookID uint) ([]entity.Chapter, int64, error) {
			return []entity.Chapter{
				{ID: 1, BookID: 1, Title: "Chapter 1"},
				{ID: 2, BookID: 1, Title: "Chapter 2"},
			}, 2, nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	result, err := uc.GetByBookID(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestChapterUsecase_GetByBookID_RepositoryError(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChaptersByBookIDFunc: func(ctx context.Context, bookID uint) ([]entity.Chapter, int64, error) {
			return nil, 0, errors.New("database error")
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	_, err := uc.GetByBookID(context.Background(), 1)

	if err == nil {
		t.Error("expected error from repository, got nil")
	}
}

// ==================== Update Tests ====================

func TestChapterUsecase_Update_Success(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{
				ID:            1,
				BookID:        1,
				ChapterNumber: 1,
				Title:         "Original Title",
				Category:      "Test",
				TotalVerses:   10,
			}, nil
		},
		UpdateChapterFunc: func(ctx context.Context, ch *entity.Chapter) error {
			ch.UpdatedAt = time.Now()
			return nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	newTitle := "Updated Title"
	input := portuc.UpdateChapterInput{
		Title: &newTitle,
	}

	result, err := uc.Update(context.Background(), 1, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", result.Title)
	}
}

func TestChapterUsecase_Update_ChapterNotFound(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return nil, errors.New("not found")
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	input := portuc.UpdateChapterInput{}

	_, err := uc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent chapter, got nil")
	}
}

func TestChapterUsecase_Update_BookNotFound(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1, BookID: 1}, nil
		},
	}
	mockBookRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			return nil, nil // Book not found
		},
	}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	newBookID := uint(999)
	input := portuc.UpdateChapterInput{
		BookID: &newBookID,
	}

	_, err := uc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for non-existent book, got nil")
	}
}

func TestChapterUsecase_Update_ValidationError(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1, BookID: 1, Title: "Test", Category: "Test", TotalVerses: 1}, nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)
	emptyTitle := ""
	input := portuc.UpdateChapterInput{
		Title: &emptyTitle, // Invalid: empty
	}

	_, err := uc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for empty title, got nil")
	}
}

// ==================== Delete Tests ====================

func TestChapterUsecase_Delete_Success(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1}, nil
		},
		DeleteChapterFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	err := uc.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestChapterUsecase_Delete_NotFound(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return nil, errors.New("not found")
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	err := uc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent chapter, got nil")
	}
}

func TestChapterUsecase_Delete_RepositoryError(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		GetChapterByIDFunc: func(ctx context.Context, id uint) (*entity.Chapter, error) {
			return &entity.Chapter{ID: 1}, nil
		},
		DeleteChapterFunc: func(ctx context.Context, id uint) error {
			return errors.New("database error")
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	err := uc.Delete(context.Background(), 1)

	if err == nil {
		t.Error("expected error from repository, got nil")
	}
}

func TestChapterUsecase_BulkDelete_Success(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		DeleteChaptersFunc: func(ctx context.Context, ids []uint) error {
			if len(ids) != 2 {
				t.Errorf("expected 2 ids, got %d", len(ids))
			}
			return nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	err := uc.BulkDelete(context.Background(), []uint{1, 2})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestChapterUsecase_BulkDelete_EmptyIDs(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		DeleteChaptersFunc: func(ctx context.Context, ids []uint) error {
			t.Error("DeleteChapters should not be called")
			return nil
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	err := uc.BulkDelete(context.Background(), []uint{})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestChapterUsecase_BulkDelete_RepositoryError(t *testing.T) {
	mockChapterRepo := &MockChapterRepository{
		DeleteChaptersFunc: func(ctx context.Context, ids []uint) error {
			return errors.New("database error")
		},
	}
	mockBookRepo := &MockBookRepository{}
	mockLogger := &MockLogger{}

	uc := chapter.NewChapterUsecase(mockChapterRepo, mockBookRepo, mockLogger)

	err := uc.BulkDelete(context.Background(), []uint{1, 2})

	if err == nil {
		t.Error("expected error from repository, got nil")
	}
}
