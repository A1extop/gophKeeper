package repository

import (
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gophKeeper/internal/client/services/lockbox/models"
	"gophKeeper/pkg/crypt"
	"log"
	"strings"
	"time"
)

var jwtSecret = []byte("your-secret-key") // Глобальный секретный ключ

type Repository interface {
	SaveLockBox(box *models.LockBox) error
	GetLockBox(name string) (*models.LockBox, error)
	Deleted(name string) error
	Updated(data *models.LockBox) error
	GetLockBoxes() (*[]models.LockBox, error)
	Exists(name string) (bool, error)
	SaveToken(token string)
	PurgeExpiredLocks() error
}

var key string = "superSecretKey19"

type SQLiteRepository struct {
	db        *sql.DB
	authToken string
	encryptor crypt.Encryptor
}

func NewSQLiteRepository(db *sql.DB) Repository {
	return &SQLiteRepository{
		db:        db,
		encryptor: crypt.New(key),
	}
}

func GetUserIDFromToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userIdFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}
	return int(userIdFloat), nil
}
func (r *SQLiteRepository) getUserID() (int, error) {
	return GetUserIDFromToken(r.authToken)
}
func (r *SQLiteRepository) SaveLockBox(box *models.LockBox) error {
	userID, err := r.getUserID()
	if err != nil {
		return err
	}
	dataEncrypt, err := crypt.EncryptLockBox(box, r.encryptor)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO lockbox (name,username, url, password, description, user_id) 
	 VALUES (?, ?, ?, ?, ?, ?)`,
		dataEncrypt.Name, dataEncrypt.Login, dataEncrypt.URL, dataEncrypt.Password, dataEncrypt.Description, userID,
	)
	return err
}

func (r *SQLiteRepository) GetLockBoxes() (*[]models.LockBox, error) {
	userID, err := r.getUserID()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(
		`SELECT id, name, username, url, password, description, created_at, updated_at 
         FROM lockbox WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lockboxes []models.LockBox
	for rows.Next() {
		var box models.LockBox
		var url, login, password, description sql.NullString

		if err := rows.Scan(
			&box.ID, &box.Name, &login, &url, &password, &description, &box.CreatedAt, &box.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if url.Valid {
			box.URL = url.String
		} else {
			box.URL = ""
		}
		if login.Valid {
			box.Login = login.String
		} else {
			box.Login = ""
		}
		if password.Valid {
			box.Password = password.String
		} else {
			box.Password = ""
		}
		if description.Valid {
			box.Description = description.String
		} else {
			box.Description = ""
		}

		dataDectypt, err := crypt.DecryptLockBox(&box, r.encryptor)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		lockboxes = append(lockboxes, *dataDectypt)
	}

	return &lockboxes, nil
}
func (r *SQLiteRepository) GetLockBox(name string) (*models.LockBox, error) {
	userID, err := r.getUserID()
	if err != nil {
		return nil, err
	}
	var box models.LockBox
	var url, username, password, description sql.NullString

	err = r.db.QueryRow(
		`SELECT id, url, username, password, description, created_at, updated_at
         FROM lockbox WHERE name = ? AND user_id = ?`, name, userID,
	).Scan(&box.ID, &url, &username, &password, &description, &box.CreatedAt, &box.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if url.Valid {
		box.URL = url.String
	} else {
		box.URL = ""
	}
	if username.Valid {
		box.Login = username.String
	} else {
		box.Login = ""
	}
	if password.Valid {
		box.Password = password.String
	} else {
		box.Password = ""
	}
	if description.Valid {
		box.Description = description.String
	} else {
		box.Description = ""
	}

	box.Name = name

	dataDecrypt, err := crypt.DecryptLockBox(&box, r.encryptor)
	if err != nil {
		return nil, err
	}

	return dataDecrypt, nil
}
func (r *SQLiteRepository) Deleted(name string) error {
	userID, err := r.getUserID()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`UPDATE lockbox SET deleted_at = ? WHERE name = ? AND user_id = ?`,
		time.Now(), name, userID,
	)
	return err
}

func (r *SQLiteRepository) Updated(data *models.LockBox) error {
	dataEncrypt, err := crypt.EncryptLockBox(data, r.encryptor)
	if err != nil {
		return err
	}
	userID, err := r.getUserID()
	if err != nil {
		return err
	}

	var query strings.Builder
	query.WriteString("UPDATE lockbox SET ")

	if data.URL != "" {
		query.WriteString(fmt.Sprintf("url = '%s', ", dataEncrypt.URL))
	}
	if data.Login != "" {
		query.WriteString(fmt.Sprintf("username = '%s', ", dataEncrypt.Login))
	}
	if data.Description != "" {
		query.WriteString(fmt.Sprintf("description = '%s', ", dataEncrypt.Description))
	}
	if data.Password != "" {
		query.WriteString(fmt.Sprintf("password = '%s', ", dataEncrypt.Password))
	}

	query.WriteString("updated_at = CURRENT_TIMESTAMP, ")
	query.WriteString("deleted_at = CASE WHEN deleted_at IS NOT NULL AND updated_at > deleted_at THEN NULL ELSE deleted_at END ")
	query.WriteString(fmt.Sprintf("WHERE name = '%s' AND user_id = %d", dataEncrypt.Name, userID))

	_, err = r.db.Exec(query.String())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *SQLiteRepository) Exists(name string) (bool, error) {
	userID, err := r.getUserID()
	if err != nil {
		return false, err
	}

	var count int
	query := "SELECT COUNT(*) FROM lockbox WHERE name = ? AND user_id = ?"
	err = r.db.QueryRow(query, name, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func (r *SQLiteRepository) SaveToken(token string) {
	r.authToken = token
}
func (r *SQLiteRepository) PurgeExpiredLocks() error {
	query := `DELETE FROM lockbox 
              WHERE deleted_at IS NOT NULL 
              AND deleted_at < NOW() - INTERVAL '1 day'`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
