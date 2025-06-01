package usecase

import (
	"context"
	"gophKeeper/internal/server/domain"
	"gophKeeper/internal/server/services/lockbox/models"
	"gophKeeper/internal/server/services/lockbox/repository"
)

type ILockBoxUsecase interface {
	UpdateLock(ctx context.Context, data *models.Data) error
	DeleteLock(ctx context.Context, name string, userId int) error
	GetLockByName(ctx context.Context, name string, userId int) (*models.Data, error)
	CreateLock(ctx context.Context, data *models.Data) (int, error)
	GetAllLocks(ctx context.Context, userId int) (*[]models.Data, error)
	ExistsLock(ctx context.Context, name string, userId int) (bool, error)
	CreateOrUpdateLock(ctx context.Context, data *models.Data) (int, error)
}

type LockBoxUsecase struct {
	repo repository.ILockBoxRepo
}

func NewLockBoxUsecase(repo repository.ILockBoxRepo) ILockBoxUsecase {
	return &LockBoxUsecase{repo: repo}
}

func (u *LockBoxUsecase) UpdateLock(ctx context.Context, data *models.Data) error {
	return u.repo.Update(ctx, data)
}
func (u *LockBoxUsecase) DeleteLock(ctx context.Context, name string, userId int) error {
	if name == "" {
		return domain.ErrNameEmpty
	}
	return u.repo.Delete(ctx, name, userId)
}
func (u *LockBoxUsecase) GetLockByName(ctx context.Context, name string, userId int) (*models.Data, error) {
	return u.repo.Get(ctx, name, userId)
}
func (u *LockBoxUsecase) CreateLock(ctx context.Context, data *models.Data) (int, error) {

	if data.Name == "" && (data.UserID == 0 || (data.Login == "" && data.Url == "" && data.Description == "" && data.Password == "")) {
		return 0, domain.ErrNoDataToCreate
	}
	return u.repo.Create(ctx, data)
}
func (u *LockBoxUsecase) GetAllLocks(ctx context.Context, userId int) (*[]models.Data, error) {
	locks, err := u.repo.GetAll(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(*locks) == 0 {
		return &[]models.Data{}, nil
	}

	return locks, nil
}

func (u *LockBoxUsecase) ExistsLock(ctx context.Context, name string, userId int) (bool, error) {
	if name == "" {
		return false, domain.ErrNameEmpty
	}
	return u.repo.Exists(ctx, name, userId)
}
func (u *LockBoxUsecase) CreateOrUpdateLock(ctx context.Context, data *models.Data) (int, error) {

	exists, err := u.ExistsLock(ctx, data.Name, data.UserID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return u.CreateLock(ctx, data)
	}

	if err := u.UpdateLock(ctx, data); err != nil {
		return 0, err
	}

	updated, err := u.GetLockByName(ctx, data.Name, data.UserID)
	if err != nil {
		return 0, err
	}
	return updated.Id, nil
}
