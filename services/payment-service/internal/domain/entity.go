package domain

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	PaymentID     string        `json:"payment_id"`
	OrderID       string        `json:"order_id"`
	Amount        float64       `json:"amount"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	PaymentDate   time.Time     `json:"payment_date"`
}

func NewPayment(orderID string, amount float64) *Payment {
	return &Payment{
		PaymentID:     uuid.New().String(),
		OrderID:       orderID,
		Amount:        amount,
		PaymentStatus: PaymentStatusPending,
		PaymentDate:   time.Now(),
	}
}

func (p *Payment) Complete() {
	p.PaymentStatus = PaymentStatusCompleted
	p.PaymentDate = time.Now()
}

func (p *Payment) Fail() {
	p.PaymentStatus = PaymentStatusFailed
	p.PaymentDate = time.Now()
}

