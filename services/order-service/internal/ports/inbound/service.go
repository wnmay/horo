package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain"
)

// OrderService defines the interface for order business logic
type OrderService interface {
	CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*domain.Order, error)
	GetOrders(ctx context.Context) ([]*domain.Order, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error)
	GetOrdersByCustomer(ctx context.Context, customerID string) ([]*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error
	UpdateOrderPaymentID(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID) error
	MarkCustomerCompleted(ctx context.Context, orderID uuid.UUID) error
	MarkProphetCompleted(ctx context.Context, orderID uuid.UUID) error
}

// CreateOrderCommand represents the command to create an order
type CreateOrderCommand struct {
	CustomerID string    `json:"customer_id" validate:"required"`
	CourseID   uuid.UUID `json:"course_id" validate:"required"`
}