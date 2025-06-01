package db

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"gophKeeper/internal/server/config"
	"log"
)

type IDatabase interface {
	GetDB() *pgxpool.Pool
	CheckMigrations(ctx context.Context, state string) error
}

type database struct {
	pool *pgxpool.Pool
}

func Init(ctx context.Context, cfg *config.Config) (IDatabase, error) {
	db, err := connectToDB(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &database{db}, err
}

func (db *database) GetDB() *pgxpool.Pool {
	return db.pool
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (db *database) CheckMigrations(ctx context.Context, state string) error {
	if err := db.pool.Ping(ctx); err != nil {
		log.Fatalf("Could not connect to PostgreSQL: %v", err)
	}

	openDBFromPool := stdlib.OpenDBFromPool(db.pool)
	fmt.Println("Connected to database")

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if state == "up" {
		if err := goose.Up(openDBFromPool, "migrations"); err != nil {
			panic(err)
		}
	} else if state == "down" {
		if err := goose.Down(openDBFromPool, "migrations"); err != nil {
			panic(err)
		}
	}

	if err := openDBFromPool.Close(); err != nil {
		panic(err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

func connectToDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s database=%s password=%s sslmode=%s",
		cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.Username, cfg.Pg.Database, cfg.Pg.Password, cfg.Pg.SSL)

	pool, err := pgxpool.New(ctx, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
