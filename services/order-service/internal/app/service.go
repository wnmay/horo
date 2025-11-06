package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain"
	"github.com/wnmay/horo/services/order-service/internal/ports/inbound"
	"github.com/wnmay/horo/services/order-service/internal/ports/outbound"
)

type OrderService struct {
	orderRepo      outbound.OrderRepository
	eventPublisher outbound.EventPublisher
	paymentService outbound.PaymentService
}

func NewOrderService(
	orderRepo outbound.OrderRepository,
	eventPublisher outbound.EventPublisher,
	paymentService outbound.PaymentService,
) inbound.OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
		paymentService: paymentService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, cmd inbound.CreateOrderCommand) (*domain.Order, error) {
	// Create new order entity
	order := domain.NewOrder(cmd.CustomerID, cmd.CourseID, cmd.RoomID)

	// Save order to repository
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Publish order created event
	if err := s.eventPublisher.PublishOrderCreated(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to publish order created event: %w", err)
	}

	// Create payment asynchronously (could be done via event as well)
	// Note: Payment amount will be determined by the payment service based on course pricing
	if err := s.paymentService.CreatePayment(ctx, order.OrderID, 0.0); err != nil {
		// Log error but don't fail the order creation
		// Payment creation will be retried via message queue
		fmt.Printf("Failed to create payment for order %s: %v\n", order.OrderID, err)
	}

	return order, nil
}
func (s *OrderService) GetOrders(ctx context.Context) ([]*domain.Order, error) {
	orders, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID string) ([]*domain.Order, error) {
	orders, err := s.orderRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for customer: %w", err)
	}
	return orders, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Update order status
	order.Status = status

	// Save updated order
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// If order is confirmed, publish order paid event
	if status == domain.StatusConfirmed {
		if err := s.eventPublisher.PublishOrderPaid(ctx, order); err != nil {
			return fmt.Errorf("failed to publish order paid event: %w", err)
		}
	}

	return nil
}

func (s *OrderService) UpdateOrderPaymentID(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID) error {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Update payment ID using the SetPaymentID method
	order.SetPaymentID(paymentID)

	// Save updated order
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order with payment ID: %w", err)
	}

	// Publish order payment bound event
	if err := s.eventPublisher.PublishOrderPaymentBound(ctx, order); err != nil {
		return fmt.Errorf("failed to publish order payment bound event: %w", err)
	}

	return nil
}

func (s *OrderService) MarkCustomerCompleted(ctx context.Context, orderID uuid.UUID) error {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Check if order is confirmed (payment completed)
	if order.Status != domain.StatusConfirmed && order.Status != domain.StatusCompleted {
		return fmt.Errorf("order must be confirmed before marking as completed")
	}

	// Mark as completed by customer
	order.MarkCustomerCompleted()
	if( order.Status == domain.StatusCompleted) {
		// Publish order completed event
		if err := s.eventPublisher.PublishOrderCompleted(ctx, order); err != nil {
			return fmt.Errorf("failed to publish order completed event: %w", err)
		}
	}

	// Save updated order
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to mark order as completed by customer: %w", err)
	}

	return nil
}

func (s *OrderService) MarkProphetCompleted(ctx context.Context, orderID uuid.UUID) error {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Check if order is confirmed (payment completed)
	if order.Status != domain.StatusConfirmed && order.Status != domain.StatusCompleted {
		return fmt.Errorf("order must be confirmed before marking as completed")
	}

	// Mark as completed by prophet
	order.MarkProphetCompleted()
	if( order.Status == domain.StatusCompleted) {
		// Publish order completed event
		if err := s.eventPublisher.PublishOrderCompleted(ctx, order); err != nil {
			return fmt.Errorf("failed to publish order completed event: %w", err)
		}
	}
	// Save updated order
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to mark order as completed by prophet: %w", err)
	}

	return nil
}