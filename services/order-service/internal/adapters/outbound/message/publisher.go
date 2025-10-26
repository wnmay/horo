package message

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wnmay/horo/services/order-service/internal/domain"
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

func (p *Publisher) PublishOrderCreated(ctx context.Context, order *domain.Order) error {
	coursePrice := 200
	// Create order data for the event
	orderData := message.OrderData{
		OrderID:    order.OrderID.String(),
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		Amount: coursePrice,
	}

	// Marshal the order data
	data, err := json.Marshal(orderData)
	if err != nil {
		return fmt.Errorf("failed to marshal order data: %w", err)
	}

	// Create AMQP message using contract structure
	amqpMessage := contract.AmqpMessage{
		OwnerID: order.CustomerID,
		Data:    data,
	}

	// Publish the message with routing key using PublishMessage method
	if err := p.rabbit.PublishMessage(ctx, contract.OrderCreatedEvent, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish order created event: %w", err)
	}

	fmt.Printf("Published order created event for order: %s\n", order.OrderID)
	return nil
}