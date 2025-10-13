package repository

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusCancelled  OrderStatus = "CANCELLED"
	StatusFailed     OrderStatus = "CONFIRMED"
)

type Order struct {
	OrderID     uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID  uuid.UUID   `gorm:"type:uuid;not null"`
	CourseID    uuid.UUID   `gorm:"type:uuid;not null"`
	PaymentID   uuid.UUID   `gorm:"type:uuid"` 
	Status      OrderStatus `gorm:"type:varchar(20);not null"`
	OrderDate   time.Time   `gorm:"not null"`
}

func (o *Order) TableName() string {
	return "ordert"
}