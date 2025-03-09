package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"gophKeeper/internal/server/services/users/models"
	"gophKeeper/util"
)

type UserUsecaseMock struct {
	mock.Mock
}

func NewUserUsecaseMock() *UserUsecaseMock {
	return &UserUsecaseMock{}
}

func (us *UserUsecaseMock) GetUsers(ctx context.Context) ([]*models.User, error) {
	args := us.Called(ctx)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (us *UserUsecaseMock) CreateUser(ctx context.Context, user *models.User) error {
	if user.UserType == "" {
		user.UserType = "attendee"
	}
	if !util.IsValidUserTypeForRegistration(user.UserType) {
		return errors.New("invalid user type")
	}
	args := us.Called(ctx, user)
	return args.Error(0)
}

func (us *UserUsecaseMock) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := us.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (us *UserUsecaseMock) UpdateUser(user *models.User) error {
	args := us.Called(user)
	return args.Error(0)
}

func (us *UserUsecaseMock) DeleteUser(id int) error {
	args := us.Called(id)
	return args.Error(0)
}
