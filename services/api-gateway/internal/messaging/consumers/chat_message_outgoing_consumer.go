package consumers

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wnmay/horo/services/api-gateway/internal/websocket"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type ChatMessageConsumer struct {
	rmq *message.RabbitMQ
	hub *websocket.Hub
}

func NewChatMessageOutgoingConsumer(rmq *message.RabbitMQ, hub *websocket.Hub) *ChatMessageConsumer {
	return &ChatMessageConsumer{
		rmq: rmq,
		hub: hub,
	}
}

func (c *ChatMessageConsumer) StartListening() error {
	return c.rmq.ConsumeMessages(message.ChatMessageOutgoingQueue, c.handleChatMessage)
}

// handleChatMessage processes outgoing chat messages and sends them to connected clients via WebSocket
func (c *ChatMessageConsumer) handleChatMessage(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("[chat-consumer] got message: %s", delivery.Body)

	var amqpMsg contract.AmqpMessage
	if err := json.Unmarshal(delivery.Body, &amqpMsg); err != nil {
		log.Printf("unmarshal amqp msg err: %v", err)
		return err
	}

	var data message.ChatMessageOutgoingData
	if err := json.Unmarshal(amqpMsg.Data, &data); err != nil {
		log.Printf("unmarshal payload err: %v", err)
		return err
	}

	payload, _ := json.Marshal(data)

	if data.RoomID != "" {
		c.hub.BroadcastToRoom(data.RoomID, payload)
	}

	return nil
}
