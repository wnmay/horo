package inbound

import (
	"context"

	"github.com/wnmay/horo/services/payment-service/internal/domain"
)

type PaymentService interface {
	CreatePaymentFromOrder(ctx context.Context, cmd CreatePaymentCommand) (*domain.Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID string, status domain.PaymentStatus) error
	CompletePayment(ctx context.Context, paymentID string) error
	SettlePayment(ctx context.Context, paymentID string) error
}

type CreatePaymentCommand struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}
