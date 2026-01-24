package usecase

import (
	"context"
	"errors"

	"ishari-backend/internal/core/entity"
	"ishari-backend/internal/core/port/repository"
	portusecase "ishari-backend/internal/core/port/usecase"
)

type bookUseCase struct {
	repo repository.BookRepository
}

func NewBookUseCase(repo repository.BookRepository) portusecase.BookUseCase {
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

func (uc *bookUseCase) CreateBook(ctx context.Context, in portusecase.CreateBookInput) (*entity.Book, error) {
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

func (uc *bookUseCase) EditBook(ctx context.Context, id int64, in portusecase.CreateBookInput) (*entity.Book, error) {
	if in.Title == "" {
		return nil, errors.New("title is required")
	}

	book := &entity.Book{
		ID:            int(id),
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
