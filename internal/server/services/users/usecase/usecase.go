package usecase

import (
	"context"
	"gophKeeper/internal/server/domain"
	"gophKeeper/internal/server/services/users/models"
	"gophKeeper/internal/server/services/users/repository"
	"gophKeeper/util"
)

type IUserUsecase interface {
	GetUsers(ctx context.Context) ([]*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

type UserUsecase struct {
	repo repository.IUserRepository
}

func NewUserUsecase(repo repository.IUserRepository) IUserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (us *UserUsecase) GetUsers(ctx context.Context) ([]*models.User, error) {
	return us.repo.GetAll(ctx)
}

func (us *UserUsecase) CreateUser(ctx context.Context, user *models.User) error {
	if user.UserType == "" {
		user.UserType = "attendee"
	}
	if !util.IsValidUserTypeForRegistration(user.UserType) {
		return domain.ErrInvalidUserType
	}
	return us.repo.Create(ctx, user)
}

func (us *UserUsecase) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return us.repo.GetByID(ctx, id)
}

func (us *UserUsecase) UpdateUser(user *models.User) error {
	return us.repo.Update(user)
}

func (us *UserUsecase) DeleteUser(id int) error {
	return us.repo.Delete(id)
}
