package config

import (
	"errors"
	"os"
)

type AppConfig struct {
	ENV         string
	PORT        string
	MONGODB_URI string
	JWT_SECRET  string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		ENV:         getEnv("ENV", "development"),
		PORT:        getEnv("PORT", "8080"),
		MONGODB_URI: getEnv("MONGODB_URI", ""),
		JWT_SECRET:  getEnv("JWT_SECRET", ""),
	}

	if cfg.MONGODB_URI == "" {
		return nil, errors.New("DB URI not found")
	}
	return cfg, nil
}

func getEnv(key, fallbackValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallbackValue
}
