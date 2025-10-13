package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	OrderID    uuid.UUID   `json:"order_id"`
	CustomerID uuid.UUID   `json:"customer_id"`
	CourseID   uuid.UUID   `json:"course_id"`
	PaymentID  *uuid.UUID  `json:"payment_id,omitempty"`
	Status     OrderStatus `json:"status"`
	OrderDate  time.Time   `json:"order_date"`
}

func NewOrder(customerID, courseID uuid.UUID) *Order {
	return &Order{
		OrderID:    uuid.New(),
		CustomerID: customerID,
		CourseID:   courseID,
		Status:     StatusPending,
		OrderDate:  time.Now(),
	}
}

func (o *Order) Confirm() {
	o.Status = StatusConfirmed
}

func (o *Order) Cancel() {
	o.Status = StatusCancelled
}

func (o *Order) SetPaymentID(paymentID uuid.UUID) {
	o.PaymentID = &paymentID
}
