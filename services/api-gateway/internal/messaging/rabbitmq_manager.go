package messaging

import (
	"fmt"
	"log"

	"github.com/wnmay/horo/services/api-gateway/internal/messaging/consumers"
	"github.com/wnmay/horo/services/api-gateway/internal/messaging/publishers"
	shared_message "github.com/wnmay/horo/shared/message"
)

// MessagingManager handles all RabbitMQ-related components including
// consumers, publishers, and the underlying RabbitMQ connection
type MessagingManager struct {
	client *shared_message.RabbitMQ

	// Consumers
	chatConsumer *consumers.ChatMessageConsumer

	// Publishers
	chatPublisher *publishers.ChatMessagePublisher
}

// NewMessagingManager creates and initializes a new messaging manager
func NewMessagingManager(rabbitmqURI string) (*MessagingManager, error) {
	// Initialize RabbitMQ client
	client, err := shared_message.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	// Initialize consumers
	chatConsumer := consumers.NewChatMessageOutgoingConsumer(client)

	// Initialize publishers
	chatPublisher := publishers.NewChatMessagePublisher(client)

	return &MessagingManager{
		client:        client,
		chatConsumer:  chatConsumer,
		chatPublisher: chatPublisher,
	}, nil
}

// StartConsumers starts all message consumers
func (m *MessagingManager) StartConsumers() error {
	log.Println("Starting message consumers...")

	if err := m.chatConsumer.StartListening(); err != nil {
		return fmt.Errorf("failed to start chat consumer: %w", err)
	}

	log.Println("All consumers started successfully")
	return nil
}

// GetChatPublisher returns the chat message publisher
// This can be used by WebSocket handlers to publish incoming messages
func (m *MessagingManager) GetChatPublisher() *publishers.ChatMessagePublisher {
	return m.chatPublisher
}

// Close gracefully shuts down all messaging components
func (m *MessagingManager) Close() error {
	log.Println("Shutting down messaging manager...")
	m.client.Close()
	return nil
}
