package usecase

import (
	"context"
	"github.com/stretchr/testify/mock"
	"gophKeeper/internal/server/services/auth/models"
)

type AuthUsecaseMock struct {
	mock.Mock
}

func NewAuthUsecaseMock() IAuthUsecase {
	return &AuthUsecaseMock{}
}

func (m *AuthUsecaseMock) CheckUser(ctx context.Context, user *models.AuthUser) (*models.InfoUser, error) {
	args := m.Called(ctx, user)
	if args.Get(0) != nil {
		return args.Get(0).(*models.InfoUser), args.Error(1)
	}
	return nil, args.Error(1)
}
