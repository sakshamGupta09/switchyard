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

const shutdownTimeout = 5 * time.Second

func main() {
	// Logging setup
	slog.SetDefault(logger.New())

	// Environment variables setup
	cfg := mustLoadConfig()

	// Database setup
	mongoClient := mustConnectMongo(&cfg.Mongo)

	// Server setup
	router := server.NewRouter()
	httpServer := server.NewServer(router, fmt.Sprint(":", cfg.Port))

	// Start Server
	serverError := startServer(httpServer)

	// Wait for server failure or shutdown signal.
	shutdownSignalCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	defer stop()

	select {
	case <-shutdownSignalCtx.Done():
		slog.Info("shutdown signal received")
		gracefulShutdown(httpServer, mongoClient)
	case err := <-serverError:
		slog.Error("server error", "error", err)
		handleServerError(err)
	}
}

func mustLoadConfig() *config.AppConfig {
	_ = godotenv.Load()
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

func startServer(httpServer *server.Server) <-chan error {
	serverError := make(chan error, 1)
	go func() {
		slog.Info("starting http server")
		serverError <- httpServer.Start()
	}()
	return serverError
}

func handleServerError(err error) {
	if errors.Is(err, http.ErrServerClosed) {
		return
	}
	slog.Error("HTTP server stopped unexpectedly", "error", err)
	os.Exit(1)
}

func gracefulShutdown(httpServer *server.Server, mongoClient *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Stop accepting new requests and wait for in-flight requests.
	if err := httpServer.Stop(ctx); err != nil {
		slog.Error("HTTP server shutdown failed", "error", err)
	}

	// Disconnect MongoDB after HTTP requests have stopped.
	if err := mongoClient.Close(ctx); err != nil {
		slog.Error("MongoDB shutdown failed", "error", err)
	}

	slog.Info("application stopped")
}
