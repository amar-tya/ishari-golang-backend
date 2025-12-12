package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/usecase"
)

// MockBookRepository is a manual mock for testing
type MockBookRepository struct {
	ListFunc    func(ctx context.Context, offset, limit int, search string) ([]entity.Book, int64, error)
	CreateFunc  func(ctx context.Context, book *entity.Book) error
	EditFunc    func(ctx context.Context, book *entity.Book) error
	DeleteFunc  func(ctx context.Context, id int64) error
	GetByIdFunc func(ctx context.Context, id int64) (*entity.Book, error)
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

func (m *MockBookRepository) GetById(ctx context.Context, id int64) (*entity.Book, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return nil, nil
}

func TestBookUseCase_ListBooks(t *testing.T) {
	mockRepo := &MockBookRepository{
		ListFunc: func(ctx context.Context, offset, limit int, search string) ([]entity.Book, int64, error) {
			return []entity.Book{
				{ID: 1, Title: "Test Book"},
			}, 1, nil
		},
	}

	uc := usecase.NewBookUseCase(mockRepo)
	books, total, err := uc.ListBooks(context.Background(), 1, 10, "")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(books) != 1 {
		t.Errorf("expected 1 book, got %d", len(books))
	}
	if books[0].Title != "Test Book" {
		t.Errorf("expected title 'Test Book', got %s", books[0].Title)
	}
}

func TestBookUseCase_CreateBook(t *testing.T) {
	mockRepo := &MockBookRepository{
		CreateFunc: func(ctx context.Context, book *entity.Book) error {
			book.ID = 1
			book.CreatedAt = time.Now()
			book.UpdatedAt = time.Now()
			return nil
		},
	}

	uc := usecase.NewBookUseCase(mockRepo)
	input := usecase.CreateBookInput{
		Title: "New Book",
	}

	book, err := uc.CreateBook(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if book.ID != 1 {
		t.Errorf("expected ID 1, got %d", book.ID)
	}
	if book.Title != "New Book" {
		t.Errorf("expected title 'New Book', got %s", book.Title)
	}
}

func TestBookUseCase_CreateBook_Error(t *testing.T) {
	mockRepo := &MockBookRepository{
		CreateFunc: func(ctx context.Context, book *entity.Book) error {
			return errors.New("db error")
		},
	}

	uc := usecase.NewBookUseCase(mockRepo)
	input := usecase.CreateBookInput{
		Title: "New Book",
	}

	_, err := uc.CreateBook(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestBookUseCase_EditBook(t *testing.T) {
	mockRepo := &MockBookRepository{
		EditFunc: func(ctx context.Context, book *entity.Book) error {
			book.UpdatedAt = time.Now()
			return nil
		},
	}

	uc := usecase.NewBookUseCase(mockRepo)
	input := usecase.CreateBookInput{
		Title: "Updated Book",
	}

	book, err := uc.EditBook(context.Background(), 1, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if book.Title != "Updated Book" {
		t.Errorf("expected title 'Updated Book', got %s", book.Title)
	}
}

func TestBookUseCase_DeleteBook(t *testing.T) {
	mockRepo := &MockBookRepository{
		DeleteFunc: func(ctx context.Context, id int64) error {
			if id != 1 {
				return errors.New("book not found")
			}
			return nil
		},
	}

	uc := usecase.NewBookUseCase(mockRepo)

	err := uc.DeleteBook(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = uc.DeleteBook(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent book, got nil")
	}
}

func TestBookUseCase_GetBookById(t *testing.T) {
	mockRepo := &MockBookRepository{
		GetByIdFunc: func(ctx context.Context, id int64) (*entity.Book, error) {
			if id == 1 {
				return &entity.Book{ID: 1, Title: "Existing Book"}, nil
			}
			return nil, errors.New("book not found")
		},
	}

	uc := usecase.NewBookUseCase(mockRepo)

	book, err := uc.GetBookById(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if book.ID != 1 {
		t.Errorf("expected ID 1, got %d", book.ID)
	}

	_, err = uc.GetBookById(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent book, got nil")
	}
}
