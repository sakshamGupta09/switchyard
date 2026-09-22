package main

import (
	"context"
	"log/slog"
	"os"
	"switchyard/internal/config"
	"switchyard/internal/database/mongo"
	"switchyard/internal/logger"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Logger setup
	appLogger := logger.New()
	slog.SetDefault(appLogger)

	// Environment variables setup
	godotenv.Load()
	config, err := config.Load()

	if err != nil {
		slog.Error("Failed to load required configurations", "error", err)
		os.Exit(1)
	}
	slog.Info("Config loaded successfully")

	// Database setup
	client, err := mongo.New(&config.Mongo)

	if err != nil {
		slog.Error("failed to connect DB", "error", err)
		os.Exit(1)
	}
	defer cleanup(client)
	slog.Info("mongoDB connected!")
}

func cleanup(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Close(ctx); err != nil {
		slog.Error("failed to close MongoDB", "error", err)
	}
}
