package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI      string
	Name     string
	Timeout  int
	PoolSize int
}

// MongoDB wraps the MongoDB client and database
type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(cfg MongoDBConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(uint64(cfg.PoolSize))

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &MongoDB{
		Client: client,
		DB:     client.Database(cfg.Name),
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// HealthCheck checks if MongoDB is healthy
func (m *MongoDB) HealthCheck(ctx context.Context) error {
	return m.Client.Ping(ctx, readpref.Primary())
}

// Collection returns a MongoDB collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.DB.Collection(name)
}
