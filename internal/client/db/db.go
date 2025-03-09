package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"log"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "mydatabase.db")
	if err != nil {
		return nil, err
	}

	migrationsDir := "internal/client/db/migrations"

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	_, err = goose.GetDBVersion(db)
	if err != nil {
		return nil, err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		if !errors.Is(err, goose.ErrNoNextVersion) {
			log.Println("Migration error:", err)
			return nil, err
		}
	}

	return db, nil
}
