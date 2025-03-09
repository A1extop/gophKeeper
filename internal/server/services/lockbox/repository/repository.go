package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gophKeeper/internal/server/db"
	"gophKeeper/internal/server/services/lockbox/models"
	"strings"
)

type ILockBoxRepo interface {
	Update(ctx context.Context, data *models.Data) error
	Delete(ctx context.Context, name string, userId int) error
	Get(ctx context.Context, name string, userId int) (*models.Data, error)
	Create(ctx context.Context, data *models.Data) (int, error)
	GetAll(ctx context.Context, userId int) (*[]models.Data, error)
	Exists(ctx context.Context, name string, userId int) (bool, error)
	PurgeExpiredLocks(ctx context.Context) (int64, error)
}

type LockBoxRepo struct {
	db db.IDatabase
}

func NewLockBoxRepo(db db.IDatabase) ILockBoxRepo {
	return &LockBoxRepo{
		db: db,
	}
}

func (l *LockBoxRepo) Update(ctx context.Context, data *models.Data) error {
	var query strings.Builder
	query.WriteString("UPDATE lockbox SET ")

	if data.Url != "" {
		query.WriteString(fmt.Sprintf("url = '%s', ", data.Url))
	}
	if data.Login != "" {
		query.WriteString(fmt.Sprintf("username = '%s', ", data.Login))
	}
	if data.Name != "" {
		query.WriteString(fmt.Sprintf("name = '%s', ", data.Name))
	}
	if data.Description != "" {
		query.WriteString(fmt.Sprintf("description = '%s', ", data.Description))
	}
	if data.Password != "" {
		query.WriteString(fmt.Sprintf("password = '%s', ", data.Password))
	}

	query.WriteString("updated_at = NOW(), ")
	query.WriteString("deleted_at = CASE WHEN deleted_at IS NOT NULL AND NOW() > deleted_at THEN NULL ELSE deleted_at END ")

	query.WriteString(fmt.Sprintf("WHERE name = '%s'", data.Name))

	_, err := l.db.GetDB().Exec(ctx, query.String())
	if err != nil {
		return err
	}

	return nil
}

// todo также нужно сверить, может быть удалено на сервере, а обновлено на локал, сверить данные, если обновление позже удаления, то это учесть
// todo а надо будет один хер как-то чистить память, возможно удалять записи по истечению какого-то времени, например раз в сутки
func (l *LockBoxRepo) Delete(ctx context.Context, name string, userId int) error {
	query := `UPDATE lockbox 
              SET deleted_at = NOW() 
              WHERE name = $1 AND user_id = $2 AND deleted_at IS NULL`

	res, err := l.db.GetDB().Exec(ctx, query, name, userId)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return nil
	}

	return nil
}

func (l *LockBoxRepo) Get(ctx context.Context, name string, userId int) (*models.Data, error) {
	query := `SELECT id, url, username, password, description, created_at, updated_at, deleted_at 
          FROM lockbox 
          WHERE user_id = $1 AND name = $2`
	var data models.Data

	err := l.db.GetDB().QueryRow(ctx, query, userId, name).Scan(&data.Id, &data.Url, &data.Login, &data.Password, &data.Description, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)
	if err != nil {
		fmt.Println(err)
		return nil, err

	}
	data.Name = name
	data.UserID = userId
	return &data, nil
}

func (l *LockBoxRepo) Create(ctx context.Context, data *models.Data) (int, error) {
	query := `INSERT INTO lockbox (name, url, username, password, description, user_id)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	err := l.db.GetDB().QueryRow(ctx, query, data.Name, data.Url, data.Login, data.Password, data.Description, data.UserID).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			restoreQuery := `UPDATE lockbox 
                             SET deleted_at = NULL, updated_at = NOW(), 
                                 url = $2, username = $3, password = $4, description = $5 
                             WHERE name = $1 AND deleted_at IS NOT NULL 
                             RETURNING id`

			err = l.db.GetDB().QueryRow(ctx, restoreQuery, data.Name, data.Url, data.Login, data.Password, data.Description).Scan(&id)
			if err == nil {
				return id, nil
			}
		}
		return 0, err
	}

	return id, nil
}

func (l *LockBoxRepo) GetAll(ctx context.Context, userId int) (*[]models.Data, error) {
	query := `SELECT id, name, url, username, password, description, created_at, updated_at, deleted_at
              FROM lockbox 
              WHERE user_id = $1 AND deleted_at IS NULL`
	rows, err := l.db.GetDB().Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []models.Data
	for rows.Next() {
		var data models.Data
		if err := rows.Scan(&data.Id, &data.Name, &data.Url, &data.Login, &data.Password, &data.Description, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt); err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &dataList, nil
}
func (l *LockBoxRepo) Exists(ctx context.Context, name string, userId int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM lockbox WHERE name = $1 AND user_id = $2)`
	var exists bool
	err := l.db.GetDB().QueryRow(ctx, query, name, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (l *LockBoxRepo) PurgeExpiredLocks(ctx context.Context) (int64, error) {
	query := `DELETE FROM lockbox 
              WHERE deleted_at IS NOT NULL 
              AND deleted_at < NOW() - INTERVAL '1 day'`

	res, err := l.db.GetDB().Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
