package outbound

import (
	"context"

	"github.com/wnmay/horo/services/payment-service/internal/domain"
)

type PaymentEventPublisher interface {
	PublishPaymentCompleted(ctx context.Context, payment *domain.Payment) error
	PublishPaymentFailed(ctx context.Context, payment *domain.Payment) error
	PublishPaymentCreated(ctx context.Context, payment *domain.Payment) error
	PublishPaymentSettled(ctx context.Context, payment *domain.Payment) error
}