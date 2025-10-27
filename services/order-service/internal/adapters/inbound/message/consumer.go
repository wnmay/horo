package message

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain"
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
	// Queue 1: Listen to payment success events (for updating order status)
	paymentSuccessQueue := message.UpdateOrderStatusQueue
	paymentSuccessRoutingKey := contract.PaymentSuccessEvent 

	// Declare the payment success queue
	if err := c.rabbit.DeclareQueue(paymentSuccessQueue, paymentSuccessRoutingKey); err != nil {
		return err
	}

	// Queue 2: Listen to payment created events (for updating payment ID)
	paymentCreatedQueue := message.UpdatePaymentIDQueue
	paymentCreatedRoutingKey := contract.PaymentCreatedEvent

	// Declare the payment created queue
	if err := c.rabbit.DeclareQueue(paymentCreatedQueue, paymentCreatedRoutingKey); err != nil {
		return err
	}

	// Start consuming payment success messages
	go func() {
		if err := c.rabbit.ConsumeMessages(paymentSuccessQueue, c.handlePaymentSuccess); err != nil {
			log.Printf("Error consuming payment success messages: %v", err)
		}
	}()

	// Start consuming payment created messages
	go func() {
		if err := c.rabbit.ConsumeMessages(paymentCreatedQueue, c.handlePaymentCreated); err != nil {
			log.Printf("Error consuming payment created messages: %v", err)
		}
	}()

	log.Println("Order service message consumers started successfully")
	return nil
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
	if err := c.orderService.UpdateOrderStatus(ctx, orderID, domain.StatusConfirmed); err != nil {
		log.Printf("Failed to update order status for order %s: %v", paymentData.OrderID, err)
		return err
	}

	log.Printf("Successfully updated order %s status to CONFIRMED after payment completion", paymentData.OrderID)
	return nil
}

func (c *Consumer) handlePaymentCreated(ctx context.Context, delivery amqp.Delivery) error {
    log.Printf("Received payment created event: %s", string(delivery.Body))

    var amqpMessage contract.AmqpMessage
    if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
        log.Printf("Failed to unmarshal AMQP message: %v", err)
        return err
    }

    var paymentData struct {
        OrderID   string `json:"orderId"`
        PaymentID string `json:"paymentId"`
    }
    if err := json.Unmarshal(amqpMessage.Data, &paymentData); err != nil {
        log.Printf("Failed to unmarshal payment data: %v", err)
        return err
    }

    log.Printf("Updating order %s with payment ID %s", paymentData.OrderID, paymentData.PaymentID)

    // Parse order ID
    orderID, err := uuid.Parse(paymentData.OrderID)
    if err != nil {
        log.Printf("Invalid order ID format: %s", paymentData.OrderID)
        return err
    }

    // Parse payment ID
    paymentID, err := uuid.Parse(paymentData.PaymentID)
    if err != nil {
        log.Printf("Invalid payment ID format: %s", paymentData.PaymentID)
        return err
    }

    // Update order with payment ID
    if err := c.orderService.UpdateOrderPaymentID(ctx, orderID, paymentID); err != nil {
        log.Printf("Failed to update order with payment ID: %v", err)
        return err
    }

    log.Printf("Successfully updated order %s with payment ID %s", paymentData.OrderID, paymentData.PaymentID)
    return nil
}