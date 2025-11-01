package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupMongoDB initializes MongoDB connection and returns a database instance
func SetupMongoDB(uri, dbName string) (*mongo.Database, *mongo.Client, error) {
	log.Println("Setting up MongoDB connection...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(context.Background())
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Successfully connected to MongoDB database: %s", dbName)

	db := client.Database(dbName)
	return db, client, nil
}
