package repository

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusCompleted PaymentStatus = "COMPLETED"
	StatusFailed    PaymentStatus = "FAILED"
	StatusCancelled PaymentStatus = "CANCELLED"
)

type Payment struct {
	PaymentID     string        `gorm:"type:varchar(255);primary_key"`
	OrderID       string        `gorm:"type:varchar(255);not null"`
	Amount        float64       `gorm:"type:decimal(10,2);not null"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(20);not null"`
	PaymentDate   time.Time     `gorm:"not null"`
}

func (p *Payment) TableName() string {
	return "payments"
}