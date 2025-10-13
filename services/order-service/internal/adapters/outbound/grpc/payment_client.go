package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/ports/outbound"
)

// PaymentClient implements the PaymentService interface
type PaymentClient struct {
	// This could be a gRPC client to payment service
	// For now, we'll simulate it
}

func NewPaymentClient() outbound.PaymentService {
	return &PaymentClient{}
}

func (p *PaymentClient) CreatePayment(ctx context.Context, orderID uuid.UUID, amount float64, customerID uuid.UUID) error {
	// TODO: Implement actual gRPC call to payment service
	// For now, just log the action
	fmt.Printf("Creating payment for order: %s, amount: %.2f, customer: %s\n", orderID, amount, customerID)
	
	// Simulate success
	return nil
}