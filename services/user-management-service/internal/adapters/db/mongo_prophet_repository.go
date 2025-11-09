package db

import (
	"context"
	"time"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProphetRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoProphetRepository(uri, dbName, collectionName string) (*MongoProphetRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoProphetRepository{
		client:     client,
		collection: collection,
	}, nil
}

// Save a Prophet into MongoDB (upsert behavior)
func (r *MongoProphetRepository) Save(ctx context.Context, prophet domain.Prophet) error {
	if prophet.User == nil {
		return mongo.ErrNilDocument
	}

	model := ProphetModel{
		User: UserModel{
			UserID:   prophet.User.ID,
			FullName: prophet.User.FullName,
			Email:    prophet.User.Email,
			Role:     prophet.User.Role,
		},
		Balance: prophet.Balance,
	}

	filter := bson.M{"user_id": model.User.UserID}
	update := bson.M{"$set": model}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindById retrieves a Prophet by ID
func (r *MongoProphetRepository) FindById(ctx context.Context, userId string) (*domain.Prophet, error) {
	var model ProphetModel

	err := r.collection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // not found
		}
		return nil, err
	}

	user := &domain.User{
		ID:       model.User.UserID,
		FullName: model.User.FullName,
		Email:    model.User.Email,
		Role:     model.User.Role,
	}

	prophet := &domain.Prophet{
		User:    user,
		Balance: model.Balance,
	}

	return prophet, nil
}
