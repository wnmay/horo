package db

import (
	"context"
	"time"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GormUserRepository → MongoUserRepository
type MongoUserRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoUserRepository connects to MongoDB and returns a repo instance.
func NewMongoUserRepository(uri, dbName, collectionName string) (*MongoUserRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoUserRepository{
		client:     client,
		collection: collection,
	}, nil
}

// Save inserts or updates a user document
func (r *MongoUserRepository) Save(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// If the document already exists, replace it; otherwise, insert it.
	filter := map[string]interface{}{"id": user.ID}
	update := map[string]interface{}{
		"$set": user,
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
