package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

// ChatPublisher handles publishing chat messages
type ChatPublisher struct {
	rmq *message.RabbitMQ
}

// NewChatPublisher creates a new chat message publisher
func NewChatPublisher(rmq *message.RabbitMQ) *ChatPublisher {
	return &ChatPublisher{
		rmq: rmq,
	}
}

// PublishOutgoingMessage publishes a message from chat service to client
func (p *ChatPublisher) PublishOutgoingMessage(ctx context.Context, msg contract.AmqpMessage) error {
	log.Printf("Publishing outgoing chat message: %+v", msg)

	err := p.rmq.PublishMessage(
		ctx,
		contract.ChatMessageOutgoingEvent,
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to publish outgoing message: %v", err)
	}

	return nil
}
