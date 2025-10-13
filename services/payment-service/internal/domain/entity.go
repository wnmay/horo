package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	ID       string        `json:"id"`
	OrderID  string        `json:"order_id"`
	UserID   string        `json:"user_id"`
	Amount   float64       `json:"amount"`
	Currency string        `json:"currency"`
	Status   string        `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

func NewPayment(orderID, userID string, amount float64, currency string) *Payment {
	return &Payment{
		ID:       uuid.New().String(),
		OrderID:  orderID,
		UserID:   userID,
		Amount:   amount,
		Currency: currency,
		Status:    string(PaymentStatusPending),
		CreatedAt: time.Now(),
	}
}

func (p *Payment) Complete() {
	p.Status = string(PaymentStatusCompleted)
	p.CreatedAt = time.Now()
}

func (p *Payment) Fail() {
	p.Status = string(PaymentStatusFailed)
	p.CreatedAt = time.Now()
}

