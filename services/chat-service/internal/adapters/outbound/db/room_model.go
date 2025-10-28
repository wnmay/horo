package repository

import "time"

type RoomModel struct {
	ID          string    `bson:"_id" gorm:"primaryKey;size:64" json:"id"` // e.g. "chat_room_123"
	ProphetID   string    `bson:"prophet_id" gorm:"size:64" json:"prophetId"`
	CustomerID  string    `bson:"customer_id" gorm:"size:64" json:"customerId"`
	CourseID    string    `bson:"course_id" gorm:"size:64" json:"courseId"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	LastMessage string    `bson:"last_message" gorm:"-" json:"lastMessage,omitempty"` // transient field, not store
}
