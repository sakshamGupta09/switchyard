package mongo

import (
	"context"
	"fmt"
	"switchyard/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	serverSelectionTimeout = 3 * time.Second
	connectTimeout         = 3 * time.Second
	pingTimeout            = 3 * time.Second
)

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func (c *Client) Close(ctx context.Context) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Disconnect(ctx)
}

func New(cfg *config.MongoConfig) (*Client, error) {
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetServerSelectionTimeout(serverSelectionTimeout).
		SetConnectTimeout(connectTimeout)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connect to MongoDB: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}
	return &Client{client: client, db: client.Database(cfg.Database)}, nil
}
