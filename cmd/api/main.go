package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"switchyard/internal/config"
	"switchyard/internal/database/mongo"
	"switchyard/internal/logger"
	"switchyard/internal/server"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	httpShutdownTimeout  = 5 * time.Second
	mongoShutdownTimeout = 5 * time.Second
)

func main() {
	slog.SetDefault(logger.New())

	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	mongoClient, err := connectMongo(&cfg.Mongo)
	if err != nil {
		return fmt.Errorf("connect MongoDB: %w", err)
	}

	router := server.NewRouter()
	httpServer := server.NewServer(router, fmt.Sprintf(":%d", cfg.Port))

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer stop()

	serverError := startServer(httpServer)

	select {
	case <-signalCtx.Done():
		slog.Info("shutdown signal received")
		gracefulShutdown(httpServer, mongoClient)
		return nil
	case err := <-serverError:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		gracefulShutdown(httpServer, mongoClient)
		return fmt.Errorf("start HTTP server: %w", err)
	}
}

func loadConfig() (*config.AppConfig, error) {
	// Load .env for local development; production uses environment variables.
	_ = godotenv.Load()
	cfg, err := config.Load()

	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func connectMongo(cfg *config.MongoConfig) (*mongo.Client, error) {
	client, err := mongo.New(cfg)

	if err != nil {
		return nil, err
	}
	return client, nil
}

func startServer(httpServer *server.Server) <-chan error {
	serverError := make(chan error, 1)
	go func() {
		slog.Info("starting http server")
		serverError <- httpServer.Start()
	}()
	return serverError
}

func gracefulShutdown(httpServer *server.Server, mongoClient *mongo.Client) {
	// Stop accepting new requests and wait for in-flight requests.
	if err := shutdownServer(httpServer); err != nil {
		slog.Error("HTTP server shutdown failed", "error", err)
	}

	// Disconnect MongoDB after HTTP requests have stopped.
	if err := disconnectMongo(mongoClient); err != nil {
		slog.Error("MongoDB shutdown failed", "error", err)
	}

	slog.Info("application stopped")
}

func shutdownServer(httpServer *server.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()

	return httpServer.Stop(ctx)
}

func disconnectMongo(mongoClient *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), mongoShutdownTimeout)
	defer cancel()

	return mongoClient.Close(ctx)
}
