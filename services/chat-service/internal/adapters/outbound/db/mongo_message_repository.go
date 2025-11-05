package repository

import (
	"context"
	"log"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
	repository_port "github.com/wnmay/horo/services/chat-service/internal/ports/outbound"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoMessageRepository struct {
	collection *mongo.Collection
}

func NewMongoMessageRepository(db *mongo.Database, collectionName string) repository_port.MessageRepository {
	collection := db.Collection(collectionName)

	return &mongoMessageRepository{
		collection: collection,
	}
}

func (r *mongoMessageRepository) SaveMessage(ctx context.Context, message *domain.Message) (string,error) {
	model := ToModel(message)
	insertResult, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return "", err
	}
	messageID := insertResult.InsertedID.(primitive.ObjectID).Hex()
	return messageID, nil
}

func (r *mongoMessageRepository) FindMessagesByRoomID(ctx context.Context, roomID string) ([]*domain.Message, error) {
	roomOID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		log.Println("Invalid RoomID for MessageModel:", err)
		roomOID = primitive.NilObjectID
	}

	filter := bson.M{"room_id": roomOID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*MessageModel
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	var domainMessages []*domain.Message
	for _, msg := range messages {
		domainMessages = append(domainMessages, ToDomain(msg))
	}

	return domainMessages, nil
}
