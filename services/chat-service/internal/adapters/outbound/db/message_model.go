package repository

import (
	"time"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoomID    string             `bson:"room_id" gorm:"index;size:64" json:"roomId"`
	SenderID  string             `bson:"sender_id" gorm:"index;size:64" json:"senderId"`
	Content   string             `bson:"content" gorm:"type:text" json:"content"`
	Type      string             `bson:"type" gorm:"size:32;default:'text'" json:"type"`     // text | notification
	Status    string             `bson:"status" gorm:"size:32;default:'sent'" json:"status"` // sent | delivered | read
	CreatedAt time.Time          `bson:"creat_at" json:"createdAt"`

	Room RoomModel `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"-"`
}

func ToMessageEntity(model MessageModel) *domain.Message {
	return &domain.Message{
		ID:        model.ID.Hex(),
		RoomID:    model.RoomID,
		SenderID:  model.SenderID,
		Content:   model.Content,
		Type:      domain.MessageType(model.Type),
		Status:    domain.MessageStatus(model.Status),
		CreatedAt: model.CreatedAt,
	}
}

func ToMessageModel(entity *domain.Message) MessageModel {
	return MessageModel{
		ID:        primitive.NewObjectID(),
		RoomID:    entity.RoomID,
		SenderID:  entity.SenderID,
		Content:   entity.Content,
		Type:      string(entity.Type),
		Status:    string(entity.Status),
		CreatedAt: entity.CreatedAt,
	}
}