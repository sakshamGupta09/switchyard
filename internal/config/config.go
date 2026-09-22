package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Environment string
	Port        string
	JWTSecret   string
	Mongo       MongoConfig
}

type MongoConfig struct {
	URI         string
	Database    string
	MinPoolSize uint64
	MaxPoolSize uint64
}

func Load() (*AppConfig, error) {
	minPoolSize, err := getUint("MONGO_MIN_POOL_SIZE")
	if err != nil {
		return nil, err
	}

	maxPoolSize, err := getUint("MONGO_MAX_POOL_SIZE")

	if err != nil {
		return nil, err
	}

	cfg := &AppConfig{
		Environment: getEnv("ENV"),
		Port:        getEnv("PORT"),
		JWTSecret:   getEnv("JWT_SECRET"),
		Mongo: MongoConfig{
			URI:         getEnv("MONGO_URI"),
			Database:    getEnv("MONGO_DATABASE"),
			MinPoolSize: minPoolSize,
			MaxPoolSize: maxPoolSize,
		},
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
		missing = append(missing, "ENV")
	}

	if cfg.Port == "" {
		missing = append(missing, "PORT")
	}

	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	if cfg.Mongo.URI == "" {
		missing = append(missing, "MONGO_URI")
	}

	if cfg.Mongo.Database == "" {
		missing = append(missing, "MONGO_DATABASE")
	}

	if len(missing) > 0 {
		return fmt.Errorf("Missing required configuration: %s",
			strings.Join(missing, ", "))
	}

	return nil
}

func getUint(key string) (uint64, error) {
	value := getEnv(key)

	if value == "" {
		return 0, fmt.Errorf("missing %s", key)
	}
	result, err := strconv.ParseUint(value, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("invalid value for %s", key)
	}
	return result, nil
}
