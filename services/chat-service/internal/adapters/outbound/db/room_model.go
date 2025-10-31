package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" gorm:"primaryKey;size:64" json:"id"` // MongoDB auto-generates if omitted
	ProphetID   string             `bson:"prophet_id" gorm:"size:64" json:"prophetId"`
	CustomerID  string             `bson:"customer_id" gorm:"size:64" json:"customerId"`
	CourseID    string             `bson:"course_id" gorm:"size:64" json:"courseId"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	LastMessage string             `bson:"last_message" gorm:"-" json:"lastMessage,omitempty"` // transient field, not store
}
