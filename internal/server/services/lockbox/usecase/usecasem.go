package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gophKeeper/internal/server/services/lockbox/models"
)

type LockBoxUsecaseMock struct {
	mock.Mock
}

// NewLockBoxUsecaseMock создает новый мок для usecase.
func NewLockBoxUsecaseMock() *LockBoxUsecaseMock {
	return &LockBoxUsecaseMock{}
}

func (u *LockBoxUsecaseMock) UpdateLock(ctx context.Context, data *models.Data) error {
	args := u.Called(ctx, data)
	return args.Error(0)
}

func (u *LockBoxUsecaseMock) DeleteLock(ctx context.Context, name string, userId int) error {
	args := u.Called(ctx, name, userId)
	return args.Error(0)
}

func (u *LockBoxUsecaseMock) GetLockByName(ctx context.Context, name string, userId int) (*models.Data, error) {
	args := u.Called(ctx, name, userId)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Data), args.Error(1)
	}
	return nil, args.Error(1)
}

func (u *LockBoxUsecaseMock) CreateLock(ctx context.Context, data *models.Data) (int, error) {
	args := u.Called(ctx, data)
	return args.Int(0), args.Error(1)
}

func (u *LockBoxUsecaseMock) GetAllLocks(ctx context.Context, userId int) (*[]models.Data, error) {
	args := u.Called(ctx, userId)
	if args.Get(0) != nil {
		return args.Get(0).(*[]models.Data), args.Error(1)
	}
	return nil, args.Error(1)
}

func (u *LockBoxUsecaseMock) ExistsLock(ctx context.Context, name string, userId int) (bool, error) {
	args := u.Called(ctx, name, userId)
	return args.Bool(0), args.Error(1)
}
