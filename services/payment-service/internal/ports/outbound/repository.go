package outbound

import (
	"context"

	"github.com/wnmay/horo/services/payment-service/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
	GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
	Delete(ctx context.Context, paymentID string) error
}