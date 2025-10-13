package message

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain/entity"
	"github.com/wnmay/horo/services/order-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type Consumer struct {
	orderService inbound.OrderService
	rabbit       *message.RabbitMQ
}

func NewConsumer(orderService inbound.OrderService, rabbit *message.RabbitMQ) *Consumer {
	return &Consumer{
		orderService: orderService,
		rabbit:       rabbit,
	}
}

func (c *Consumer) StartListening() error {
	// Listen to the UpdateOrderStatusQueue for payment success events
	queueName := message.UpdateOrderStatusQueue

	// Declare the queue first before consuming
	if err := c.rabbit.DeclareQueue(queueName, queueName); err != nil {
		return err
	}

	// Start consuming messages
	return c.rabbit.ConsumeMessages(queueName, c.handlePaymentSuccess)
}

func (c *Consumer) handlePaymentSuccess(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("Received payment success event: %s", delivery.Body)

	var amqpMessage contract.AmqpMessage
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Printf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	var paymentData message.PaymentSuccessData
	if err := json.Unmarshal(amqpMessage.Data, &paymentData); err != nil {
		log.Printf("Failed to unmarshal payment success data: %v", err)
		return err
	}

	log.Printf("Processing payment success event for order: %s, transaction: %s", 
		paymentData.OrderID, paymentData.TransactionID)

	// Parse OrderID
	orderID, err := uuid.Parse(paymentData.OrderID)
	if err != nil {
		log.Printf("Invalid order ID format: %s", paymentData.OrderID)
		return err
	}

	// Update order status to confirmed
	if err := c.orderService.UpdateOrderStatus(ctx, orderID, entity.StatusConfirmed); err != nil {
		log.Printf("Failed to update order status for order %s: %v", paymentData.OrderID, err)
		return err
	}

	log.Printf("Successfully updated order %s status to CONFIRMED", paymentData.OrderID)
	return nil
}