package outbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain/entity"
)

// OrderRepository defines the interface for order data persistence
type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, orderID uuid.UUID) (*entity.Order, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, orderID uuid.UUID) error
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, order *entity.Order) error
}

// PaymentService defines the interface for payment operations
type PaymentService interface {
	CreatePayment(ctx context.Context, orderID uuid.UUID, amount float64, customerID uuid.UUID) error
}