package app

import "go.mongodb.org/mongo-driver/v2/mongo"

type RouteDeps struct {
	MongoClient *mongo.Client
}
