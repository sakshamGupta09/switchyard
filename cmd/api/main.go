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

const shutdownTimeout = 5 * time.Second

func main() {
	// Logging setup
	slog.SetDefault(logger.New())

	// Environment variables setup
	cfg := mustLoadConfig()

	// Database setup
	client := mustConnectMongo(&cfg.Mongo)
	defer closeMongo(client)
}

func mustLoadConfig() *config.AppConfig {
	godotenv.Load()
	cfg, err := config.Load()

	if err != nil {
		slog.Error("loading config", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded successfully")
	return cfg
}

func mustConnectMongo(cfg *config.MongoConfig) *mongo.Client {
	client, err := mongo.New(cfg)

	if err != nil {
		slog.Error("connecting to MongoDB", "error", err)
		os.Exit(1)
	}
	slog.Info("MongoDB connected!")
	return client
}

func closeMongo(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := client.Close(ctx); err != nil {
		slog.Error("disconnecting MongoDB", "error", err)
		return
	}
}
