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

	var typeCheck struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(amqpMsg.Data, &typeCheck); err != nil {
		log.Printf("unmarshal type check err: %v", err)
		return err
	}

	var payload []byte
	var roomID string

	switch typeCheck.Type {
	case "text":
		// Handle regular text messages
		var data message.ChatMessageOutgoingData
		if err := json.Unmarshal(amqpMsg.Data, &data); err != nil {
			log.Printf("unmarshal text message err: %v", err)
			return err
		}
		roomID = data.RoomID
		payload, _ = json.Marshal(data)
		log.Printf("[chat-consumer] text message: roomID=%s, senderID=%s", data.RoomID, data.SenderID)

	case "notification":
		// Handle notification messages with detailed data
		var data map[string]interface{}
		if err := json.Unmarshal(amqpMsg.Data, &data); err != nil {
			log.Printf("unmarshal notification message err: %v", err)
			return err
		}
		if roomIDVal, ok := data["roomId"].(string); ok {
			roomID = roomIDVal
		}
		payload, _ = json.Marshal(data)
		log.Printf("[chat-consumer] notification: roomID=%s, trigger=%v", roomID, data["trigger"])

	default:
		log.Printf("unknown message type: %s", typeCheck.Type)
		return nil
	}

	if roomID != "" {
		c.hub.BroadcastToRoom(roomID, payload)
	}

	return nil
}
