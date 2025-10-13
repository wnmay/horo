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

func (r *MongoUserRepository) Save(ctx context.Context, user domain.User) error {
	filter := map[string]interface{}{"id": user.ID}
	update := map[string]interface{}{
		"$set": user,
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
