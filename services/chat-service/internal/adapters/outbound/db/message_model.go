package repository

import (
	"time"

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
