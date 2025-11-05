package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string
const (
	StatusPending   OrderStatus = "PENDING"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusCompleted  OrderStatus = "COMPLETED"
)

type Order struct {
	OrderID              uuid.UUID   `json:"order_id"`
	CustomerID           string      `json:"customer_id"`
	CourseID             string       `json:"course_id"`
	PaymentID            *uuid.UUID  `json:"payment_id,omitempty"`
	Status               OrderStatus `json:"status"`
	IsCustomerCompleted  bool        `json:"is_customer_completed"`
	IsProphetCompleted   bool        `json:"is_prophet_completed"`
	CustomerCompletedAt  *time.Time  `json:"customer_completed_at,omitempty"`
	ProphetCompletedAt   *time.Time  `json:"prophet_completed_at,omitempty"`
	OrderDate            time.Time   `json:"order_date"`
}

func NewOrder(customerID string, courseID string) *Order {
	return &Order{
		OrderID:             uuid.New(),
		CustomerID:          customerID,
		CourseID:            courseID,
		Status:              StatusPending,
		IsCustomerCompleted: false,
		IsProphetCompleted:  false,
		OrderDate:           time.Now(),
	}
}

func (o *Order) Confirm() {
	o.Status = StatusConfirmed
}

func (o *Order) Cancel() {
	o.Status = StatusCancelled
}

func (o *Order) MarkCustomerCompleted() {
	now := time.Now()
	o.IsCustomerCompleted = true
	o.CustomerCompletedAt = &now
	o.checkAndMarkComplete()
}

func (o *Order) MarkProphetCompleted() {
	now := time.Now()
	o.IsProphetCompleted = true
	o.ProphetCompletedAt = &now
	o.checkAndMarkComplete()
}

func (o *Order) checkAndMarkComplete() {
	// Only mark as completed if both prophet and customer have completed
	if o.IsCustomerCompleted && o.IsProphetCompleted {
		o.Status = StatusCompleted
	}	
}

func (o *Order) Complete() {
	o.Status = StatusCompleted
}

func (o *Order) SetPaymentID(paymentID uuid.UUID) {
	o.PaymentID = &paymentID
}
