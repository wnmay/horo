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
	createPaymentQueue := message.CreatePaymentQueue
	OrderCreateRoutingKey := contract.OrderCreatedEvent

	// Declare the queue
	if err := c.rabbit.DeclareQueue(createPaymentQueue, OrderCreateRoutingKey); err != nil {
		return err
	}
    if err := c.rabbit.ConsumeMessages(createPaymentQueue, c.handleOrderCreated); err != nil {
        return err
    }

	settlePaymentQueue := message.SettlePaymentQueue
	OrderCompleteRoutingKey := contract.OrderCompletedEvent

	if err := c.rabbit.DeclareQueue(settlePaymentQueue, OrderCompleteRoutingKey); err != nil {
		return err
	}
    if err := c.rabbit.ConsumeMessages(settlePaymentQueue, c.handleOrderCompleted); err != nil {
        return err
    }
 
	return nil

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

func (c *Consumer) handleOrderCompleted(ctx context.Context, delivery amqp.Delivery) error {
    log.Printf("Received order completed event: %s", delivery.Body)

    var amqpMessage contract.AmqpMessage
    if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
        log.Printf("Failed to unmarshal AMQP message: %v", err)
        return err
    }

    var orderCompletedData message.OrderCompletedData
    if err := json.Unmarshal(amqpMessage.Data, &orderCompletedData); err != nil {
        log.Printf("Failed to unmarshal order data: %v", err)
        return err
    }

    log.Printf("Processing order completed event for order: %s", orderCompletedData.OrderID)

    if err := c.paymentService.SettlePayment(ctx, orderCompletedData.OrderID, orderCompletedData.ProphetID); err != nil {
        log.Printf("Failed to complete payment %s for order %s: %v",
            orderCompletedData.OrderID, orderCompletedData.OrderID, err)
        return err
    }

    log.Printf("Successfully completed payment for order %s", orderCompletedData.OrderID)
    return nil
}