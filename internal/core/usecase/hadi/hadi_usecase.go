package hadi

import (
	"context"
	"errors"

	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/adapter/handler/http/response"
	"ishari-backend/internal/core/entity"
	portrepo "ishari-backend/internal/core/port/repository"
	portusecase "ishari-backend/internal/core/port/usecase"

	"gorm.io/gorm"
)

type hadiUseCase struct {
	hadiRepo portrepo.HadiRepository
}

// NewHadiUseCase creates a new hadi usecase instance
func NewHadiUseCase(hadiRepo portrepo.HadiRepository) portusecase.HadiUseCase {
	return &hadiUseCase{
		hadiRepo: hadiRepo,
	}
}

func (u *hadiUseCase) Create(ctx context.Context, req dto.CreateHadiRequest) (*response.HadiResponse, error) {
	hadi := &entity.Hadi{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
	}

	err := u.hadiRepo.Create(ctx, hadi)
	if err != nil {
		return nil, err
	}

	return response.MapHadiResponse(hadi), nil
}

func (u *hadiUseCase) GetByID(ctx context.Context, id int) (*response.HadiResponse, error) {
	hadi, err := u.hadiRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hadi not found")
		}
		return nil, err
	}

	return response.MapHadiResponse(hadi), nil
}

func (u *hadiUseCase) List(ctx context.Context, page, limit int) ([]response.HadiResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	hadis, total, err := u.hadiRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return response.MapHadiListResponse(hadis), total, nil
}

func (u *hadiUseCase) Update(ctx context.Context, id int, req dto.UpdateHadiRequest) (*response.HadiResponse, error) {
	hadi, err := u.hadiRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hadi not found")
		}
		return nil, err
	}

	hadi.Name = req.Name
	hadi.Description = req.Description
	hadi.ImageURL = req.ImageURL

	err = u.hadiRepo.Update(ctx, hadi)
	if err != nil {
		return nil, err
	}

	return response.MapHadiResponse(hadi), nil
}

func (u *hadiUseCase) Delete(ctx context.Context, id int) error {
	_, err := u.hadiRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("hadi not found")
		}
		return err
	}

	err = u.hadiRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
