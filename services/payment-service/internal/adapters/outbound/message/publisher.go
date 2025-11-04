package message

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
	sharedMessage "github.com/wnmay/horo/shared/message"
)

type Publisher struct {
	rabbit *sharedMessage.RabbitMQ
}

func NewPublisher(rabbit *sharedMessage.RabbitMQ) *Publisher {
	return &Publisher{
		rabbit: rabbit,
	}
}

func (p *Publisher) PublishPaymentCompleted(ctx context.Context, payment *domain.Payment) error {
	// Create payment completion data
	paymentData := message.PaymentPublishedData{
		PaymentID:  payment.PaymentID,
		OrderID:    payment.OrderID,
		Status:     string(payment.Status),
		Amount:     payment.Amount,
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

func (p *Publisher) PublishPaymentFailed(ctx context.Context, payment *domain.Payment) error {
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

func (p *Publisher) PublishPaymentCreated(ctx context.Context, payment *domain.Payment) error {
	paymentData := message.PaymentPublishedData{
		PaymentID:  payment.PaymentID,
		OrderID:    payment.OrderID,
		Status:     string(payment.Status),
		Amount:     payment.Amount,
	}
	data, err := json.Marshal(paymentData)
	if err != nil {
		return fmt.Errorf("failed to marshal payment data: %w", err)
	}

	amqpMessage := contract.AmqpMessage{
		OwnerID: payment.OrderID,
		Data:    data,
	}

	if err := p.rabbit.PublishMessage(ctx, contract.PaymentCreatedEvent, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish payment created event: %w", err)
	}

	return nil
}

func (p *Publisher) PublishPaymentSettled(ctx context.Context, payment *domain.Payment) error {
	payload := message.PaymentPublishedData{
		PaymentID:  payment.PaymentID,
		OrderID:    payment.OrderID,
		Status:     string(payment.Status),
		Amount:     payment.Amount,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payment settled payload: %w", err)
	}

	msg := contract.AmqpMessage{
		OwnerID: payment.OrderID,
		Data:    data,
	}

	if err := p.rabbit.PublishMessage(ctx, contract.PaymentSettledEvent, msg); err != nil {
		return fmt.Errorf("failed to publish payment settled event: %w", err)
	}

	fmt.Printf("Published payment settled event for order: %s, payment: %s\n",
		payment.OrderID, payment.PaymentID)
	return nil
}
