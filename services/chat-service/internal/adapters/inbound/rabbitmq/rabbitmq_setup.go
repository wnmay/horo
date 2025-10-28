package rabbitmq

import (
	"fmt"
	"log"

	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

// SetupChatQueues sets up all queues specific to the chat service
func SetupChatQueues(rmq *message.RabbitMQ) error {
	log.Println("Setting up chat service queues...")

	// Setup incoming message queue (binds to AppExchange by default)
	if err := rmq.DeclareQueue(
		message.ChatMessageIncomingQueue,
		contract.ChatMessageIncomingEvent,
	); err != nil {
		return fmt.Errorf("failed to setup incoming queue: %v", err)
	}

	// Setup outgoing message queue (binds to AppExchange by default)
	if err := rmq.DeclareQueue(
		message.ChatMessageOutgoingQueue,
		contract.ChatMessageOutgoingEvent,
	); err != nil {
		return fmt.Errorf("failed to setup outgoing queue: %v", err)
	}

	return nil
}
