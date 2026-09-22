package mongo

import (
	"context"
	"switchyard/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		SetServerSelectionTimeout(5 * time.Second).
		SetConnectTimeout(5 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		client.Disconnect(shutdownCtx)
		return nil, err
	}
	return &Client{client: client, db: client.Database(cfg.Database)}, nil

}
