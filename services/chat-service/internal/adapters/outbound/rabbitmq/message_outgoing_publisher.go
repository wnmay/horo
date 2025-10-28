package rabbitmq

import (
	"context"
	"fmt"
	"log"

	outbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/outbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

// ChatPublisher handles publishing chat messages
type ChatPublisher struct {
	rmq *message.RabbitMQ
}

// NewChatPublisher creates a new chat message publisher
func NewChatPublisher(rmq *message.RabbitMQ) outbound_port.MessagePublisher {
	return &ChatPublisher{
		rmq: rmq,
	}
}

// PublishOutgoingMessage publishes a message from chat service to client
func (p *ChatPublisher) Publish(ctx context.Context, message contract.AmqpMessage) error {
	log.Printf("Publishing outgoing chat message: %+v", message)

	err := p.rmq.PublishMessage(
		ctx,
		contract.ChatMessageOutgoingEvent,
		message,
	)
	if err != nil {
		return fmt.Errorf("failed to publish outgoing message: %v", err)
	}

	return nil
}
