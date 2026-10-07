package health

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	mongoClient *mongo.Client
}

func NewService(mongoClient *mongo.Client) *Service {
	return &Service{
		mongoClient: mongoClient,
	}
}

func (s *Service) CheckMongoConnection(ctx context.Context) error {
	return s.mongoClient.Ping(ctx, nil)
}
