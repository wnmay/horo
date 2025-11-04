package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusSettled PaymentStatus = "SETTLED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	PaymentID  string        `json:"payment_id"`
	OrderID    string        `json:"order_id"`
	ProphetID  string        `json:"prophet_id"`
	CourseID   string        `json:"course_id"`
	CustomerID string        `json:"customer_id"`
	Amount     float64       `json:"amount"`
	Status     PaymentStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}


var (
	ErrInvalidTransition = errors.New("invalid payment status transition")
)

func NewPayment(orderID string, amount float64) *Payment {
	now := time.Now()
	return &Payment{
		PaymentID: uuid.New().String(),
		OrderID:   orderID,
		Amount:    amount,
		Status:    PaymentStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *Payment) Complete() error {
	if p.Status == PaymentStatusCompleted {
		return nil
	}
	if p.Status != PaymentStatusPending {
		return ErrInvalidTransition
	}
	p.Status = PaymentStatusCompleted
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) Settle() error {
	if p.Status == PaymentStatusSettled {
		return nil
	}
	if p.Status != PaymentStatusCompleted {
		return ErrInvalidTransition
	}
	p.Status = PaymentStatusSettled
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) Fail() {
	p.Status = PaymentStatusFailed
	p.UpdatedAt = time.Now()
}

