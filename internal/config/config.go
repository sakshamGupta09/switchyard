package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Environment string
	Port        uint64
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
	var validationErrors []error

	environment, err := requiredString("ENV")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	jwtSecret, err := requiredString("JWT_SECRET")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	mongoURI, err := requiredString("MONGO_URI")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	dbName, err := requiredString("MONGO_DATABASE")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	minPoolSize, err := requiredUint("MONGO_MIN_POOL_SIZE")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	maxPoolSize, err := requiredUint("MONGO_MAX_POOL_SIZE")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	port, err := requiredUint("PORT")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("invalid configuration: %w", errors.Join(validationErrors...))
	}

	cfg := &AppConfig{
		Environment: environment,
		JWTSecret:   jwtSecret,
		Port:        port,
		Mongo: MongoConfig{
			URI:         mongoURI,
			Database:    dbName,
			MinPoolSize: minPoolSize,
			MaxPoolSize: maxPoolSize,
		},
	}

	if err = validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func validate(cfg *AppConfig) error {
	var errs []error

	if cfg.Port < 1 || cfg.Port > 65535 {
		errs = append(errs,
			fmt.Errorf("PORT must be between 1 and 65535"),
		)
	}

	if len(cfg.JWTSecret) < 32 {
		errs = append(errs,
			fmt.Errorf("JWT secret must contain at least 32 characters"),
		)
	}

	if cfg.Mongo.MinPoolSize > cfg.Mongo.MaxPoolSize {
		errs = append(errs,
			fmt.Errorf(
				"MONGO_MIN_POOL_SIZE cannot exceed MONGO_MAX_POOL_SIZE",
			),
		)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func requiredString(key string) (string, error) {
	value := getEnv(key)
	if value == "" {
		return "", fmt.Errorf("missing required variable %s", key)
	}
	return value, nil
}

func requiredUint(key string) (uint64, error) {
	value := getEnv(key)

	if value == "" {
		return 0, fmt.Errorf("missing required variable %s", key)
	}
	intValue, err := strconv.ParseUint(value, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}
	return intValue, nil
}

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
