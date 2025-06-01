package repository

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"gophKeeper/internal/server/db"
	"gophKeeper/internal/server/domain"
	"gophKeeper/internal/server/services/auth/models"
)

type IAuthRepo interface {
	GetInfoUser(ctx context.Context, user *models.AuthUser) (*models.InfoUser, error)
}

type authRepository struct {
	db db.IDatabase
}

func NewAuthRepository(db db.IDatabase) IAuthRepo {
	return &authRepository{
		db: db,
	}
}

func (a *authRepository) GetInfoUser(ctx context.Context, user *models.AuthUser) (*models.InfoUser, error) {
	var infoUser models.InfoUser
	var storedPasswordHash string

	query := `SELECT user_id, username, user_type, password_hash FROM "users" WHERE username = $1`
	row := a.db.GetDB().QueryRow(ctx, query, user.Username)
	err := row.Scan(&infoUser.UserId, &infoUser.Username, &infoUser.UserType, &storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(user.Password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return &infoUser, nil
}
