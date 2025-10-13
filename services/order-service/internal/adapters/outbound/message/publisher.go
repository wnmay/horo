package message

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wnmay/horo/services/order-service/internal/domain/entity"
	"github.com/wnmay/horo/services/order-service/internal/ports/outbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type Publisher struct {
	rabbit *message.RabbitMQ
}

func NewPublisher(rabbit *message.RabbitMQ) outbound.EventPublisher {
	return &Publisher{
		rabbit: rabbit,
	}
}

func (p *Publisher) PublishOrderCreated(ctx context.Context, order *entity.Order) error {
	// Create order data for the event
	orderData := message.OrderData{
		OrderID:    order.OrderID.String(),
		CustomerID: order.CustomerID.String(),
		Amount:     order.Amount,
		Status:     string(order.Status),
	}

	// Marshal the order data
	data, err := json.Marshal(orderData)
	if err != nil {
		return fmt.Errorf("failed to marshal order data: %w", err)
	}

	// Create AMQP message using contract structure
	amqpMessage := contract.AmqpMessage{
		OwnerID: order.CustomerID.String(),
		Data:    data,
	}

	// Publish the message with routing key using PublishMessage method
	if err := p.rabbit.PublishMessage(ctx, contract.OrderCreatedEvent, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish order created event: %w", err)
	}

	fmt.Printf("Published order created event for order: %s\n", order.OrderID)
	return nil
}

func (p *Publisher) PublishOrderStatusChanged(ctx context.Context, order *entity.Order) error {
	// Create order data for the event
	orderData := message.OrderData{
		OrderID:    order.OrderID.String(),
		CustomerID: order.CustomerID.String(),
		Amount:     order.Amount,
		Status:     string(order.Status),
	}

	// Marshal the order data
	data, err := json.Marshal(orderData)
	if err != nil {
		return fmt.Errorf("failed to marshal order data: %w", err)
	}

	// Create AMQP message using contract structure
	amqpMessage := contract.AmqpMessage{
		OwnerID: order.CustomerID.String(),
		Data:    data,
	}

	// Publish the message (you might want a different routing key for status changes)
	routingKey := fmt.Sprintf("order.status.%s", order.Status)
	if err := p.rabbit.PublishMessage(ctx, routingKey, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish order status changed event: %w", err)
	}

	fmt.Printf("Published order status changed event for order: %s, status: %s\n", order.OrderID, order.Status)
	return nil
}