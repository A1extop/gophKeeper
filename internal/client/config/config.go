package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port   string
	PgHost string
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
func New() *Config {
	_ = godotenv.Load(".env.client")
	return &Config{
		Port:   getEnv("HTTP_PORT", "8090"),
		PgHost: getEnv("PG_HOST", "http://localhost"),
	}
}
