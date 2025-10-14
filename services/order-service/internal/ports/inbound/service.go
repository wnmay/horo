package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain/entity"
)

// OrderService defines the interface for order business logic
type OrderService interface {
	CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*entity.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*entity.Order, error)
	GetOrdersByCustomer(ctx context.Context, customerID string) ([]*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entity.OrderStatus) error
}

// CreateOrderCommand represents the command to create an order
type CreateOrderCommand struct {
	CustomerID string    `json:"customer_id" validate:"required"`
	CourseID   uuid.UUID `json:"course_id" validate:"required"`
}