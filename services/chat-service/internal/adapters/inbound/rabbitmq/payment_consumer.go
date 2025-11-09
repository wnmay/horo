package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type paymentConsumer struct {
	chatService inbound_port.ChatService
	rmq         *message.RabbitMQ
}

func NewPaymentConsumer(chatService inbound_port.ChatService, rmq *message.RabbitMQ) inbound_port.MessageConsumer {
	return &paymentConsumer{
		chatService: chatService,
		rmq:         rmq,
	}
}

func (c *paymentConsumer) StartListening() error {
	return c.rmq.ConsumeMessages(message.NotifyCreatePayment, c.handlePaymentCreated)
}

func (c *paymentConsumer) handlePaymentCreated(ctx context.Context, delivery amqp.Delivery) error {
	var amqpMessage contract.AmqpMessage
	var paymentData message.PaymentPublishedData

	// Parse the AMQP message
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Fatalf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	if err := json.Unmarshal(amqpMessage.Data, &paymentData); err != nil {
		log.Printf("Failed to unmarshal message data: %v", err)
		return err
	}

	err := c.chatService.PublishPaymentCreatedMessage(ctx, paymentData.PaymentID, paymentData.OrderID, paymentData.Status, paymentData.Amount)
	if err != nil {
		log.Printf("Failed to publish payment created message: %v", err)
		return err
	}

	log.Printf("Published payment created message: %s", paymentData.PaymentID)
	return nil
}
