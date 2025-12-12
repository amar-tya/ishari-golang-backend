package usecase

import (
	"context"
	"errors"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
)

type BookUseCase interface {
	ListBooks(ctx context.Context, page, limit int, search string) ([]entity.Book, int64, error)
	CreateBook(ctx context.Context, in CreateBookInput) (*entity.Book, error)
	EditBook(ctx context.Context, id int64, in CreateBookInput) (*entity.Book, error)
	DeleteBook(ctx context.Context, id int64) error
	GetBookById(ctx context.Context, id int64) (*entity.Book, error)
}

type bookUseCase struct {
	repo repository.BookRepository
}

func NewBookUseCase(repo repository.BookRepository) BookUseCase {
	return &bookUseCase{repo: repo}
}

func (uc *bookUseCase) ListBooks(ctx context.Context, page, limit int, search string) ([]entity.Book, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return uc.repo.List(ctx, offset, limit, search)
}

type CreateBookInput struct {
	Title         string
	Author        *string
	Description   *string
	PublishedYear *int
	CoverImageURL *string
}

func (uc *bookUseCase) CreateBook(ctx context.Context, in CreateBookInput) (*entity.Book, error) {
	if in.Title == "" {
		return nil, errors.New("title is required")
	}

	book := &entity.Book{
		Title:         in.Title,
		Author:        in.Author,
		Description:   in.Description,
		PublishedYear: in.PublishedYear,
		CoverImageURL: in.CoverImageURL,
	}

	if err := uc.repo.Create(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (uc *bookUseCase) EditBook(ctx context.Context, id int64, in CreateBookInput) (*entity.Book, error) {
	if in.Title == "" {
		return nil, errors.New("title is required")
	}

	book := &entity.Book{
		Title:         in.Title,
		Author:        in.Author,
		Description:   in.Description,
		PublishedYear: in.PublishedYear,
		CoverImageURL: in.CoverImageURL,
	}

	if err := uc.repo.Edit(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (uc *bookUseCase) DeleteBook(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *bookUseCase) GetBookById(ctx context.Context, id int64) (*entity.Book, error) {
	return uc.repo.GetById(ctx, id)
}
