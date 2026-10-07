package app

import "go.mongodb.org/mongo-driver/v2/mongo"

type Dependencies struct {
	MongoClient *mongo.Client
}
