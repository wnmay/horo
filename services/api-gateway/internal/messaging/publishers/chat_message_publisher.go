package publishers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

// ChatMessagePublisher handles publishing chat messages to the message broker
type ChatMessagePublisher struct {
	rmq *message.RabbitMQ
}

// NewChatMessagePublisher creates a new chat message publisher
func NewChatMessagePublisher(rmq *message.RabbitMQ) *ChatMessagePublisher {
	return &ChatMessagePublisher{
		rmq: rmq,
	}
}

// PublishMessageIncoming publishes an incoming chat message to the chat service
// This will be called by WebSocket handlers when a client sends a message
func (p *ChatMessagePublisher) PublishMessageIncoming(
	ctx context.Context,
	roomID string,
	senderID string,
	content string,
	messageType string,
) error {
	// Create the message data
	messageData := message.ChatMessageIncomingData{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
		Type:     messageType,
	}

	// Marshal the message data
	data, err := json.Marshal(messageData)
	if err != nil {
		return fmt.Errorf("failed to marshal message data: %v", err)
	}

	// Create the AMQP message envelope
	amqpMessage := contract.AmqpMessage{
		Data: data,
	}

	log.Printf("Publishing incoming chat message - Room: %s, Sender: %s, Type: %s",
		roomID, senderID, messageType)

	// Publish to the message broker
	err = p.rmq.PublishMessage(ctx, contract.ChatMessageIncomingEvent, amqpMessage)
	if err != nil {
		return fmt.Errorf("failed to publish incoming message: %v", err)
	}

	log.Printf("Successfully published incoming message from sender %s to room %s", senderID, roomID)
	return nil
}
