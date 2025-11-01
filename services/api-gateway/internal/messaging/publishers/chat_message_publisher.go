package publishers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type ChatMessagePublisher struct {
	rmq *message.RabbitMQ
}

func NewChatMessagePublisher(rmq *message.RabbitMQ) *ChatMessagePublisher {
	return &ChatMessagePublisher{
		rmq: rmq,
	}
}

// Called by ws, publish message to chat service
func (p *ChatMessagePublisher) PublishMessageIncoming(
	ctx context.Context,
	roomID string,
	senderID string,
	content string,
	messageType string,
) error {
	messageData := message.ChatMessageIncomingData{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
		Type:     messageType,
	}

	data, err := json.Marshal(messageData)
	if err != nil {
		return fmt.Errorf("failed to marshal message data: %v", err)
	}

	amqpMessage := contract.AmqpMessage{
		Data: data,
	}

	log.Printf("[DEBUG] Publishing incoming chat message - Room: %s, Sender: %s, Type: %s",
		roomID, senderID, messageType)

	// Publish to the message broker
	err = p.rmq.PublishMessage(ctx, contract.ChatMessageIncomingEvent, amqpMessage)
	if err != nil {
		return fmt.Errorf("failed to publish incoming message: %v", err)
	}
	return nil
}
