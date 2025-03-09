package usecase

import (
	"context"
	"gophKeeper/internal/server/services/auth/models"
	"gophKeeper/internal/server/services/auth/repository"
)

type IAuthUsecase interface {
	CheckUser(ctx context.Context, user *models.AuthUser) (*models.InfoUser, error)
}

type AuthUsecase struct {
	repo repository.IAuthRepo
}

func NewAuthUsecase(repo repository.IAuthRepo) IAuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}

func (a *AuthUsecase) CheckUser(ctx context.Context, user *models.AuthUser) (*models.InfoUser, error) {
	infoUser, err := a.repo.GetInfoUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return infoUser, nil
}
