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

type Consumer struct {
	chatService inbound_port.ChatService
	rmq         *message.RabbitMQ
}

func NewMessageIncomingConsumer(chatService inbound_port.ChatService, rmq *message.RabbitMQ) inbound_port.MessageConsumer {
	return &Consumer{
		chatService: chatService,
		rmq:         rmq,
	}
}

func (c *Consumer) StartListening() error {
	return c.rmq.ConsumeMessages(message.ChatMessageIncomingQueue, c.handleMessageIncoming)
}

func (c *Consumer) handleMessageIncoming(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("Received message incoming: %s", delivery.Body)
	var amqpMessage contract.AmqpMessage
	var messageIncoming message.ChatMessageIncomingData

	// Parse the AMQP message
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Fatalf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	if err := json.Unmarshal(amqpMessage.Data, &messageIncoming); err != nil {
		log.Printf("Failed to unmarshal message data: %v", err)
		return err
	}

	roomID := messageIncoming.RoomID
	senderID := messageIncoming.SenderID
	content := messageIncoming.Content

	// Update order status to confirmed
	if err := c.chatService.SaveMessage(ctx, roomID, senderID, content); err != nil {
		log.Printf("Failed to save message %s: %v", content, err)
		return err
	}

	return nil
}
