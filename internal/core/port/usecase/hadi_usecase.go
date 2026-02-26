package usecase

import (
	"context"

	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/adapter/handler/http/response"
)

type HadiUseCase interface {
	Create(ctx context.Context, req dto.CreateHadiRequest) (*response.HadiResponse, error)
	GetByID(ctx context.Context, id int) (*response.HadiResponse, error)
	List(ctx context.Context, page, limit int) ([]response.HadiResponse, int64, error)
	Update(ctx context.Context, id int, req dto.UpdateHadiRequest) (*response.HadiResponse, error)
	Delete(ctx context.Context, id int) error
}
