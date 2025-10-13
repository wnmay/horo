package message

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wnmay/horo/services/payment-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type Consumer struct {
	paymentService inbound.PaymentService
	rabbit         *message.RabbitMQ
}

func NewConsumer(paymentService inbound.PaymentService, rabbit *message.RabbitMQ) *Consumer {
	return &Consumer{
		paymentService: paymentService,
		rabbit:         rabbit,
	}
}

func (c *Consumer) StartListening() error {
	// The queue setup should be done in the RabbitMQ setup
	// For now, we'll directly consume from a predefined queue
	queueName := message.CreatePaymentQueue

	// Start consuming messages
	return c.rabbit.ConsumeMessages(queueName, c.handleOrderCreated)
}

func (c *Consumer) handleOrderCreated(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("Received order created event: %s", delivery.Body)

	var amqpMessage contract.AmqpMessage
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Printf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	var orderData message.OrderData
	if err := json.Unmarshal(amqpMessage.Data, &orderData); err != nil {
		log.Printf("Failed to unmarshal order data: %v", err)
		return err
	}

	log.Printf("Processing order created event for order: %s, amount: %.2f", 
		orderData.OrderID, orderData.Amount)

	// Create payment command
	cmd := inbound.CreatePaymentCommand{
		OrderID: orderData.OrderID,
		Amount:  orderData.Amount,
	}

	// Call payment service to create payment
	payment, err := c.paymentService.CreatePaymentFromOrder(ctx, cmd)
	if err != nil {
		log.Printf("Failed to create payment for order %s: %v", orderData.OrderID, err)
		return err
	}

	log.Printf("Successfully created payment %s for order %s", payment.PaymentID, orderData.OrderID)
	return nil
}