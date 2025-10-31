package db

import (
	"context"
	"log"
	"time"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GormUserRepository → MongoUserRepository
type MongoUserRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserRepository(uri, dbName, collectionName string) (*MongoUserRepository, error) {
	log.Println("[DEBUG]", uri)
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

// in your repository layer
func (r *MongoUserRepository) Save(ctx context.Context, user domain.User) error {
	model := UserModel{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}

	filter := bson.M{"user_id": model.UserID}
	update := bson.M{"$set": model}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
