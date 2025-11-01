package consumers

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type ChatMessageConsumer struct {
	rmq *message.RabbitMQ
	// TODO: Add websocket connection manager here when implemented
	// wsManager *ws_connection.ConnectionManager
}

func NewChatMessageOutgoingConsumer(rmq *message.RabbitMQ) *ChatMessageConsumer {
	return &ChatMessageConsumer{
		rmq: rmq,
	}
}

func (c *ChatMessageConsumer) StartListening() error {
	return c.rmq.ConsumeMessages(message.ChatMessageOutgoingQueue, c.handleChatMessage)
}

// handleChatMessage processes outgoing chat messages and sends them to connected clients via WebSocket
func (c *ChatMessageConsumer) handleChatMessage(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("[DEBUG] Received chat message: %s", delivery.Body)

	var amqpMessage contract.AmqpMessage
	var messageData message.ChatMessageOutgoingData

	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Printf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	if err := json.Unmarshal(amqpMessage.Data, &messageData); err != nil {
		log.Printf("Failed to unmarshal message data: %v", err)
		return err
	}
	// TODO: Send message to connected WebSocket clients
	// This is where you would use the WebSocket connection manager to send
	// the message to the appropriate client(s)
	// Example:
	// if c.wsManager != nil {
	//     c.wsManager.SendToUser(messageData.RecipientID, messageData)
	// }

	return nil
}
