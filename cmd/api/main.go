package main

import (
	"log"
	"log/slog"
	"switchyard/internal/config"
	"switchyard/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config, err := config.Load()

	if err != nil {
		log.Fatal("Failed to load required configurations", err)
	}

	logger := logger.New()
	slog.SetDefault(logger)

	slog.Info("Config loaded successfully")
}
