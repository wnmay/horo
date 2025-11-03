package repository

import (
	"log"
	"time"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    primitive.ObjectID `bson:"room_id,omitempty" json:"roomId"`
	SenderID  string             `bson:"sender_id" json:"senderId"`
	Content   string             `bson:"content" json:"content"`
	Type      string             `bson:"type" json:"type"`     // text | notification
	Status    string             `bson:"status" json:"status"` // sent | delivered | read
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

func ToDomain(model *MessageModel) *domain.Message {
	return &domain.Message{
		ID:        model.ID.Hex(),
		RoomID:    model.RoomID.Hex(),
		SenderID:  model.SenderID,
		Content:   model.Content,
		Type:      domain.MessageType(model.Type),
		Status:    domain.MessageStatus(model.Status),
		CreatedAt: model.CreatedAt,
	}
}

func ToModel(entity *domain.Message) *MessageModel {
	roomOID, err := primitive.ObjectIDFromHex(entity.RoomID)
	if err != nil {
		log.Println("Invalid RoomID for MessageModel:", err)
		roomOID = primitive.NilObjectID
	}
	return &MessageModel{
		ID:        primitive.NewObjectID(),
		RoomID:    roomOID,
		SenderID:  entity.SenderID,
		Content:   entity.Content,
		Type:      string(entity.Type),
		Status:    string(entity.Status),
		CreatedAt: entity.CreatedAt,
	}
}
