package messaging

import (
	"fmt"
	"log"

	"github.com/wnmay/horo/services/api-gateway/internal/messaging/consumers"
	"github.com/wnmay/horo/services/api-gateway/internal/messaging/publishers"
	"github.com/wnmay/horo/services/api-gateway/internal/websocket"
	shared_message "github.com/wnmay/horo/shared/message"
)

// MessagingManager handles all RabbitMQ-related components including
// consumers, publishers, and the underlying RabbitMQ connection
type MessagingManager struct {
    client        *shared_message.RabbitMQ
    chatPublisher *publishers.ChatMessagePublisher
}

func NewMessagingManager(rabbitmqURI string) (*MessagingManager, error) {
    client, err := shared_message.NewRabbitMQ(rabbitmqURI)
    if err != nil {
        return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
    }

    chatPublisher := publishers.NewChatMessagePublisher(client)

    return &MessagingManager{
        client:        client,
        chatPublisher: chatPublisher,
    }, nil
}

func (m *MessagingManager) GetChatPublisher() *publishers.ChatMessagePublisher {
    return m.chatPublisher
}
// Close gracefully shuts down all messaging components
func (m *MessagingManager) Close() error {
	if m.client == nil {
        return nil
    }
	log.Println("Shutting down messaging manager...")
	m.client.Close()
	return nil
}

func (m *MessagingManager) RabbitMQ() *shared_message.RabbitMQ {
    return m.client
}

func (m *MessagingManager) StartChatConsumer(hub *websocket.Hub) {
    chatConsumer := consumers.NewChatMessageOutgoingConsumer(m.RabbitMQ(), hub)
    go func() {
        if err := chatConsumer.StartListening(); err != nil {
            log.Fatalf("chat consumer failed: %v", err)
        }
    }()
}
