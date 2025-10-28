package message

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wnmay/horo/services/chat-service/internal/domain"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type Consumer struct {
	chatService inbound_port.ChatService
	rabbit      *message.RabbitMQ
}

type messageIncomingData struct {
	RoomID   string `json:"roomId"`
	SenderID string `json:"senderId"`
	Content  string `json:"content"`
	Type     string `json:"type"` // text | notification
}

func NewConsumer(chatService inbound_port.ChatService, rabbit *message.RabbitMQ) *inbound_port.MessageConsumer {
	return &Consumer{
		chatService: chatService,
		rabbit:      rabbit,
	}
}

func (c *Consumer) StartListening() error {
	// Listen to chat incoming
	chatMessageIncomingQueue := message.ChatMessageIncomingQueue
	chatMessageIncomingEvent := contract.ChatMessageIncomingEvent

	// Listen to order created events (for sending notifications message)
	// paymentCreatedQueue := message.UpdatePaymentIDQueue
	// paymentCreatedRoutingKey := contract.PaymentCreatedEvent

	// Start consuming payment success messages
	go func() {
		if err := c.rabbit.ConsumeMessages(chatMessageIncomingQueue, c.handleMessageIncoming); err != nil {
			log.Printf("Error consuming payment success messages: %v", err)
		}
	}()

	// Start consuming payment created messages
	// go func() {
	// 	if err := c.rabbit.ConsumeMessages(paymentCreatedQueue, c.handlePaymentCreated); err != nil {
	// 		log.Printf("Error consuming payment created messages: %v", err)
	// 	}
	// }()

	log.Println("Order service message consumers started successfully")
	return nil
}

func (c *Consumer) handleMessageIncoming(ctx context.Context, delivery amqp.Delivery) error {
	var amqpMessage contract.AmqpMessage
	var messageIncoming messageIncomingData

	// Parse the AMQP message
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Fatalf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	if err := json.Unmarshal(amqpMessage.Data, &messageIncoming); err != nil {
		log.Printf("Failed to unmarshal payment completion data: %v", err)
		return err
	}

	roomID := messageIncoming.RoomID
	senderID := messageIncoming.SenderID
	content := messageIncoming.Content

	// Update order status to confirmed
	if err := c.chatService.SaveMessage(ctx, roomID, senderID, content); err != nil {
		log.Printf("Failed to update order status for order %s: %v", messageIncomingData.content, err)
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
