package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type EventPublisher struct {
	rabbit *message.RabbitMQ
}

func NewEventPublisher(rabbit *message.RabbitMQ) *EventPublisher {
	return &EventPublisher{
		rabbit: rabbit,
	}
}

func (p *EventPublisher) PublishPaymentCompleted(ctx context.Context, payment *domain.Payment) error {
	// Create payment success data
	paymentData := message.PaymentSuccessData{
		OrderID:       payment.OrderID,
		PaymentMethod: "credit_card",
		TransactionID: payment.PaymentID, // Using PaymentID as TransactionID for now
	}

	// Marshal the payment data
	data, err := json.Marshal(paymentData)
	if err != nil {
		return fmt.Errorf("failed to marshal payment data: %w", err)
	}

	// Create AMQP message using contract structure
	amqpMessage := contract.AmqpMessage{
		OwnerID: payment.OrderID, // Use OrderID as owner
		Data:    data,
	}

	// Publish the message with routing key using PublishMessage method
	if err := p.rabbit.PublishMessage(ctx, contract.PaymentSuccessEvent, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish payment success event: %w", err)
	}

	fmt.Printf("Published payment success event for order: %s, payment: %s\n", payment.OrderID, payment.PaymentID)
	return nil
}

func (p *EventPublisher) PublishPaymentFailed(ctx context.Context, payment *domain.Payment) error {
	// Create payment failure data (using a simple structure for now)
	paymentFailureData := map[string]interface{}{
		"order_id":   payment.OrderID,
		"payment_id": payment.PaymentID,
		"reason":     "payment_processing_failed",
	}

	// Marshal the payment failure data
	data, err := json.Marshal(paymentFailureData)
	if err != nil {
		return fmt.Errorf("failed to marshal payment failure data: %w", err)
	}

	// Create AMQP message using contract structure
	amqpMessage := contract.AmqpMessage{
		OwnerID: payment.OrderID, // Use OrderID as owner
		Data:    data,
	}

	// Publish the message with routing key for payment failure
	if err := p.rabbit.PublishMessage(ctx, "payment.failed", amqpMessage); err != nil {
		return fmt.Errorf("failed to publish payment failed event: %w", err)
	}

	fmt.Printf("Published payment failed event for order: %s, payment: %s\n", payment.OrderID, payment.PaymentID)
	return nil
}