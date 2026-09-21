package config

import (
	"fmt"
	"os"
	"strings"
)

type AppConfig struct {
	Environment string
	Port        string
	MongoDBURI  string
	JWTSecret   string
}

func Load() (*AppConfig, error) {
	cfg := &AppConfig{
		Environment: getEnv("ENV"),
		Port:        getEnv("PORT"),
		MongoDBURI:  getEnv("MONGODB_URI"),
		JWTSecret:   getEnv("JWT_SECRET"),
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func validate(cfg *AppConfig) error {
	var missing []string

	if cfg.Environment == "" {
		missing = append(missing, "Environment")
	}

	if cfg.Port == "" {
		missing = append(missing, "Port")
	}

	if cfg.MongoDBURI == "" {
		missing = append(missing, "MongoDBURI")
	}

	if cfg.JWTSecret == "" {
		missing = append(missing, "JWTSecret")
	}

	if len(missing) > 0 {
		return fmt.Errorf("Missing required configuration: %s",
			strings.Join(missing, ", "))
	}

	return nil
}
