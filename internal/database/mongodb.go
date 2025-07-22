package database

import (
	"context"
	"go-microservice-boilerplate/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(config config.MongoDBConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	database := client.Database(config.Database)

	return &MongoDB{
		Client:   client,
		Database: database,
	}, nil
}

func (m *MongoDB) Disconnect() error {
	return m.Client.Disconnect(context.Background())
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
