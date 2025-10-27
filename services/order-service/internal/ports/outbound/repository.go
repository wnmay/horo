package outbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain"
)

// OrderRepository defines the interface for order data persistence
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetAll(ctx context.Context) ([]*domain.Order, error)
	GetByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, orderID uuid.UUID) error
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, order *domain.Order) error
}

// PaymentService defines the interface for payment operations
type PaymentService interface {
	CreatePayment(ctx context.Context, orderID uuid.UUID, amount float64) error
}