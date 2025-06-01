package usecase

import (
	"context"
	"errors"
	errors1 "gophKeeper/internal/client/errors"
	"gophKeeper/internal/client/services/lockbox/clients"
	"gophKeeper/internal/client/services/lockbox/models"
	"gophKeeper/internal/client/services/lockbox/repository"
	"time"

	"log"
)

type ILockBoxUsecase interface {
	CreateLockBox(ctx context.Context, data *models.LockBoxInput) (int, error)
	DeleteLockBox(ctx context.Context, name string) error
	GetLockBoxById(ctx context.Context, name string) (*models.LockBox, error)
	GetLockBoxAll(ctx context.Context) (*[]models.LockBox, error)
	UpdateLockBox(ctx context.Context, data *models.LockBoxInput) error
	Register(ctx context.Context, username, password string) error
	Authenticate(ctx context.Context, username, password string) error
	IsAuthenticated() bool
	SyncUpdatesToServer(ctx context.Context) error
	SyncUpdatesToLocal(ctx context.Context) error
}
type LockboxUsecase struct {
	lockBoxService    clients.LockBoxService
	lockBoxRepository repository.Repository
}

func NewLockboxUsecase(lockBoxService clients.LockBoxService, lockBoxRepository repository.Repository) ILockBoxUsecase {
	return &LockboxUsecase{lockBoxService: lockBoxService, lockBoxRepository: lockBoxRepository}
}

func (uc *LockboxUsecase) CreateLockBox(ctx context.Context, data *models.LockBoxInput) (int, error) {
	if data.Name == "" {
		return 0, errors1.ErrNameLockboxRequired
	}
	if data.URL == "" && data.Login == "" && data.Password == "" && data.Description == "" {
		return 0, errors1.ErrDataRequired
	}
	id, err := uc.lockBoxService.Create(ctx, data)
	if err != nil {
		if errors.Is(err, errors1.ErrLockboxNameTakenByUser) {
			return 0, errors1.ErrExists
		}

		lockBox := models.LockBox{
			Name:        data.Name,
			Login:       data.Login,
			URL:         data.URL,
			Password:    data.Password,
			Description: data.Description,
		}

		err = uc.lockBoxRepository.SaveLockBox(&lockBox)

		if err != nil {
			log.Println("failed to save lockbox locally:", err)
			return 0, err
		}

		return lockBox.ID, nil
	}

	lockBox := models.LockBox{
		ID:          id,
		Name:        data.Name,
		URL:         data.URL,
		Password:    data.Password,
		Description: data.Description,
		SyncedAt:    time.Now(),
	}

	err = uc.lockBoxRepository.SaveLockBox(&lockBox)
	if err != nil {
		log.Println("failed to save lockbox locally:", err)
	}

	return id, nil
}

func (uc *LockboxUsecase) DeleteLockBox(ctx context.Context, name string) error {
	err := uc.lockBoxRepository.Deleted(name)
	if err != nil {
		log.Println(err)
	}
	return uc.lockBoxService.Delete(ctx, name)
}

func (uc *LockboxUsecase) GetLockBoxById(ctx context.Context, name string) (*models.LockBox, error) {
	lockBox, err1 := uc.lockBoxService.Get(ctx, name)
	if err1 != nil {
		log.Println(err1)
		lockBox, err2 := uc.lockBoxRepository.GetLockBox(name)
		if err2 != nil {
			log.Println(err2)
			return nil, errors1.ErrNotFound
		}
		return lockBox, nil
	}
	return lockBox, nil
}

func (uc *LockboxUsecase) GetLockBoxAll(ctx context.Context) (*[]models.LockBox, error) {
	lockBoxes, err1 := uc.lockBoxService.GetAll(ctx)
	if err1 != nil {
		log.Println(err1)
		lockBoxes, err2 := uc.lockBoxRepository.GetLockBoxes()
		if err2 != nil {
			log.Println(err2)
			return nil, errors1.ErrNotFound
		}
		return lockBoxes, nil
	}
	return lockBoxes, nil

}

func (uc *LockboxUsecase) UpdateLockBox(ctx context.Context, data *models.LockBoxInput) error {
	if data.Name == "" || (data.Login == "" && data.Password == "" && data.Description == "" && data.URL == "") {
		return errors1.ErrNodataToUpdate
	}
	lockBox := models.LockBox{
		Name:        data.Name,
		Login:       data.Login,
		URL:         data.URL,
		Password:    data.Password,
		Description: data.Description,
	}
	err := uc.lockBoxService.Update(ctx, data)
	if err != nil {
		log.Println("failed to update lockbox service:", err)
		err1 := uc.lockBoxRepository.Updated(&lockBox)
		if err1 != nil {
			log.Println("failed to update lockbox local:", err1)
			return errors1.ErrNotFound
		}
		return nil
	}
	lockBox.SyncedAt = time.Now()
	_ = uc.lockBoxRepository.SaveLockBox(&lockBox)
	return nil
}

func (uc *LockboxUsecase) Register(ctx context.Context, username, password string) error {
	if len(username) < 3 {
		return errors1.ErrUsernameTooShort
	}
	if len(password) < 6 {
		return errors1.ErrPasswordTooShort
	}

	return uc.lockBoxService.RegisterUser(ctx, username, password)
}

func (uc *LockboxUsecase) Authenticate(ctx context.Context, username, password string) error {
	if len(username) < 3 {
		return errors1.ErrIncorrectUsername
	}
	if len(password) < 6 {
		return errors1.ErrIncorrectPassword
	}
	token, err := uc.lockBoxService.AuthUser(ctx, username, password)
	if err != nil {
		return err
	}
	uc.lockBoxRepository.SaveToken(token)
	return nil
}

func (uc *LockboxUsecase) IsAuthenticated() bool {
	return uc.lockBoxService.Authenticated()
}

func (uc *LockboxUsecase) SyncUpdatesToServer(ctx context.Context) error {
	items, err := uc.lockBoxRepository.GetLockBoxes()
	if err != nil {
		return err
	}
	for _, item := range *items {
		err = uc.lockBoxService.UpdateOrCreate(ctx, &item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *LockboxUsecase) SyncUpdatesToLocal(ctx context.Context) error {
	items, err := uc.lockBoxService.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, item := range *items {
		exists, err := uc.lockBoxRepository.Exists(item.Name)
		if err != nil {
			return err
		}
		if !exists {
			err = uc.lockBoxRepository.SaveLockBox(&item)
			if err != nil {
				return err
			}
		}
		err = uc.lockBoxRepository.Updated(&item)
		if err != nil {
			return err
		}
	}
	return nil
}
