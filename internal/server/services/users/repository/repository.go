package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gophKeeper/internal/server/db"
	"gophKeeper/internal/server/services/users/models"
)

type IUserRepository interface {
	GetAll(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userId int) (*models.User, error)
	Update(user *models.User) error
	Delete(userId int) error
}

type userRepository struct {
	db db.IDatabase
}

func NewUserRepository(db db.IDatabase) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	query := `
        INSERT INTO users (username, password_hash, user_type)
        VALUES ($1,  $2, $3) RETURNING user_id`

	return ur.db.GetDB().QueryRow(ctx, query, user.Username, user.Password, user.UserType).Scan(&user.UserId)
}

func (ur *userRepository) GetByID(ctx context.Context, userId int) (*models.User, error) {
	var user models.User
	query := `SELECT user_id, username, email, user_type FROM "users" WHERE user_id = $1`
	err := ur.db.GetDB().QueryRow(ctx, query, userId).Scan(&user.UserId, &user.Username, &user.UserType)

	if err != nil {
		if err == sql.ErrNoRows {
			// Если нет строк, это означает, что логин и/или пароль неверны
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) Update(user *models.User) error {
	//UPDATE users SET field1 = @f1, field2 = @f2
	//WHERE Id = @Id
	return nil
}

func (ur *userRepository) Delete(id int) error {
	//DELETE FROM users WHERE Id = @Id
	return nil
}

func (ur *userRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User

	query := `
		SELECT user_id, username, password_hash, user_type, created_at
		FROM "users"
	`

	rows, err := ur.db.GetDB().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UserId,
			&user.Username,
			&user.Password,
			&user.UserType,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
