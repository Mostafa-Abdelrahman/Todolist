package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	Environment string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "todos.db"),
		Environment: getEnv("ENVIRONMENT", "development"),
		}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
