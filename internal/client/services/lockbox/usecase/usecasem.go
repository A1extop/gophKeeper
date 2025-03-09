package usecase

import (
	"context"
	"fmt"
	"gophKeeper/internal/client/services/lockbox/models"
	"time"
)

type MockLockBoxUsecase struct{}

func NewLockBoxUsecaseMock() ILockBoxUsecase {
	return &MockLockBoxUsecase{}
}
func (m *MockLockBoxUsecase) CreateLockBox(ctx context.Context, data *models.LockBoxInput) (int, error) {
	return 123, nil
}

func (m *MockLockBoxUsecase) DeleteLockBox(ctx context.Context, name string) error {
	if name == "error" {
		return fmt.Errorf("delete error")
	}
	return nil
}

func (m *MockLockBoxUsecase) GetLockBoxById(ctx context.Context, name string) (*models.LockBox, error) {
	if name == "notfound" {
		return nil, fmt.Errorf("not found")
	}
	return &models.LockBox{
		Name:        name,
		URL:         "http://example.com",
		Login:       "user",
		Password:    "pass",
		Description: "test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (m *MockLockBoxUsecase) GetLockBoxAll(ctx context.Context) (*[]models.LockBox, error) {
	boxes := []models.LockBox{
		{
			Name:        "box1",
			URL:         "http://example.com",
			Login:       "user1",
			Password:    "pass1",
			Description: "desc1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "box2",
			URL:         "http://example.org",
			Login:       "user2",
			Password:    "pass2",
			Description: "desc2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	return &boxes, nil
}

func (m *MockLockBoxUsecase) UpdateLockBox(ctx context.Context, data *models.LockBoxInput) error {
	return nil
}

func (m *MockLockBoxUsecase) Register(ctx context.Context, username, password string) error {
	return nil
}

func (m *MockLockBoxUsecase) Authenticate(ctx context.Context, username, password string) error {
	return nil
}

func (m *MockLockBoxUsecase) IsAuthenticated() bool {
	return true
}

func (m *MockLockBoxUsecase) SyncUpdatesToServer(ctx context.Context) error {
	return nil
}

func (m *MockLockBoxUsecase) SyncUpdatesToLocal(ctx context.Context) error {
	return nil
}
