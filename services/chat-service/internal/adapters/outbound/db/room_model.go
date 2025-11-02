package repository

import (
	"time"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
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

func (r *RoomModel) ToDomain() *domain.Room {
	return &domain.Room{
		ID:         r.ID.Hex(),
		ProphetID:  r.ProphetID,
		CustomerID: r.CustomerID,
		CourseID:   r.CourseID,
		CreatedAt:  r.CreatedAt,
	}
}
