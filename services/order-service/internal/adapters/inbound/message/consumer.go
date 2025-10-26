package message

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain/entity"
	"github.com/wnmay/horo/services/order-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/message"
	"github.com/wnmay/horo/shared/contract"
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
	routingKey := contract.PaymentSuccessEvent 

	// Declare the queue first before consuming
	if err := c.rabbit.DeclareQueue(queueName, routingKey); err != nil {
		return err
	}

	// Start consuming messages
	return c.rabbit.ConsumeMessages(queueName, c.handlePaymentSuccess)
}

func (c *Consumer) handlePaymentSuccess(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("Received payment completion event: %s", delivery.Body)

	// Parse the AMQP message
	var amqpMessage struct {
		OwnerID string `json:"ownerId"`
		Data    []byte `json:"data"`
	}
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Printf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	// Parse payment completion data
	var paymentData struct {
		PaymentID string  `json:"payment_id"`
		OrderID   string  `json:"order_id"`
		Status    string  `json:"status"`
		Amount    float64 `json:"amount"`
	}

	if err := json.Unmarshal(amqpMessage.Data, &paymentData); err != nil {
		log.Printf("Failed to unmarshal payment completion data: %v", err)
		return err
	}

	log.Printf("Processing payment completion for order: %s, payment: %s, status: %s", 
		paymentData.OrderID, paymentData.PaymentID, paymentData.Status)

	// Only process if payment is completed
	if paymentData.Status != "COMPLETED" {
		log.Printf("Payment status is %s, not processing", paymentData.Status)
		return nil
	}

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

	log.Printf("Successfully updated order %s status to CONFIRMED after payment completion", paymentData.OrderID)
	return nil
}