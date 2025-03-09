package config

import (
	"os"

	"github.com/joho/godotenv"
)

type PgConf struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSL      string
	DSN      string
	Migrate  string
}

type AppConf struct {
	Host                string
	Port                string
	Mode                string
	SignaturePrivateKey string
	SignaturePublicKey  string
	LogLevel            string
}

type Config struct {
	Pg  PgConf
	App AppConf
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func New() *Config {
	_ = godotenv.Load(".env.server")
	return &Config{
		Pg: PgConf{
			Host:     getEnv("PG_HOST", "postgres"),
			Port:     getEnv("PG_PORT", "5432"),
			Username: getEnv("PG_USERNAME", "user"),
			Password: getEnv("PG_PASSWORD", "password"),
			Database: getEnv("PG_DATABASE", "events_db"),
			SSL:      getEnv("PG_SSL", "disable"),
			Migrate:  getEnv("PG_MIGRATE", "up"), // up, down
		},
		App: AppConf{
			Host:     getEnv("HTTP_HOST", "0.0.0.0"),
			Port:     getEnv("HTTP_PORT", "8080"),
			Mode:     getEnv("APP_MODE", "debug"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
	}
}
