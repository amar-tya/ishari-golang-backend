package hadi_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ishari-backend/internal/adapter/handler/http/dto"
	"ishari-backend/internal/core/entity"
	hadiusecase "ishari-backend/internal/core/usecase/hadi"
)

// MockHadiRepository is a manual mock for testing
type MockHadiRepository struct {
	CreateFunc  func(ctx context.Context, hadi *entity.Hadi) error
	GetByIDFunc func(ctx context.Context, id int) (*entity.Hadi, error)
	ListFunc    func(ctx context.Context, limit, offset int) ([]entity.Hadi, int64, error)
	UpdateFunc  func(ctx context.Context, hadi *entity.Hadi) error
	DeleteFunc  func(ctx context.Context, id int) error
}

func (m *MockHadiRepository) Create(ctx context.Context, hadi *entity.Hadi) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, hadi)
	}
	return nil
}

func (m *MockHadiRepository) GetByID(ctx context.Context, id int) (*entity.Hadi, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockHadiRepository) List(ctx context.Context, limit, offset int) ([]entity.Hadi, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, 0, nil
}

func (m *MockHadiRepository) Update(ctx context.Context, hadi *entity.Hadi) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, hadi)
	}
	return nil
}

func (m *MockHadiRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestHadiUseCase_Create(t *testing.T) {
	mockRepo := &MockHadiRepository{
		CreateFunc: func(ctx context.Context, hadi *entity.Hadi) error {
			hadi.ID = 1
			hadi.CreatedAt = time.Now()
			hadi.UpdatedAt = time.Now()
			return nil
		},
	}

	uc := hadiusecase.NewHadiUseCase(mockRepo)
	desc := "A renowned reciter"
	imgURL := "http://example.com/image.png"

	req := dto.CreateHadiRequest{
		Name:        "Test Hadi",
		Description: &desc,
		ImageURL:    &imgURL,
	}

	resp, err := uc.Create(context.Background(), req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
	if resp.Name != "Test Hadi" {
		t.Errorf("expected Name 'Test Hadi', got %s", resp.Name)
	}
	if resp.Description == nil || *resp.Description != desc {
		t.Errorf("expected Description %s, got %v", desc, resp.Description)
	}
}

func TestHadiUseCase_GetByID(t *testing.T) {
	mockRepo := &MockHadiRepository{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.Hadi, error) {
			if id == 1 {
				return &entity.Hadi{ID: 1, Name: "Existing Hadi"}, nil
			}
			return nil, errors.New("record not found")
		},
	}

	uc := hadiusecase.NewHadiUseCase(mockRepo)

	resp, err := uc.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}

	_, err = uc.GetByID(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent hadi, got nil")
	}
}

func TestHadiUseCase_List(t *testing.T) {
	mockRepo := &MockHadiRepository{
		ListFunc: func(ctx context.Context, limit, offset int) ([]entity.Hadi, int64, error) {
			return []entity.Hadi{
				{ID: 1, Name: "Test Hadi 1"},
				{ID: 2, Name: "Test Hadi 2"},
			}, 2, nil
		},
	}

	uc := hadiusecase.NewHadiUseCase(mockRepo)

	resp, total, err := uc.List(context.Background(), 1, 10)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(resp) != 2 {
		t.Errorf("expected 2 hadis, got %d", len(resp))
	}
	if resp[0].Name != "Test Hadi 1" {
		t.Errorf("expected name 'Test Hadi 1', got %s", resp[0].Name)
	}
}

func TestHadiUseCase_Update(t *testing.T) {
	mockRepo := &MockHadiRepository{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.Hadi, error) {
			if id == 1 {
				return &entity.Hadi{ID: 1, Name: "Old Name"}, nil
			}
			return nil, errors.New("record not found")
		},
		UpdateFunc: func(ctx context.Context, hadi *entity.Hadi) error {
			hadi.UpdatedAt = time.Now()
			return nil
		},
	}

	uc := hadiusecase.NewHadiUseCase(mockRepo)
	desc := "Updated Description"

	req := dto.UpdateHadiRequest{
		Name:        "Updated Name",
		Description: &desc,
	}

	resp, err := uc.Update(context.Background(), 1, req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp.Name != "Updated Name" {
		t.Errorf("expected Name 'Updated Name', got %s", resp.Name)
	}
	if resp.Description == nil || *resp.Description != desc {
		t.Errorf("expected Description %s, got %v", desc, resp.Description)
	}
}

func TestHadiUseCase_Delete(t *testing.T) {
	mockRepo := &MockHadiRepository{
		GetByIDFunc: func(ctx context.Context, id int) (*entity.Hadi, error) {
			if id == 1 {
				return &entity.Hadi{ID: 1}, nil
			}
			return nil, errors.New("record not found")
		},
		DeleteFunc: func(ctx context.Context, id int) error {
			if id != 1 {
				return errors.New("hadi not found")
			}
			return nil
		},
	}

	uc := hadiusecase.NewHadiUseCase(mockRepo)

	err := uc.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = uc.Delete(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent hadi, got nil")
	}
}
