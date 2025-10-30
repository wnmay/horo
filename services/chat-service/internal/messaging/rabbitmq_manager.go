package messaging

import (
	"fmt"
	"log"

	consumerRabbit "github.com/wnmay/horo/services/chat-service/internal/adapters/inbound/rabbitmq"
	publisherRabbit "github.com/wnmay/horo/services/chat-service/internal/adapters/outbound/rabbitmq"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	outbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/outbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

// MessagingManager handles all RabbitMQ-related components including
type MessagingManager struct {
	client *message.RabbitMQ

	// Consumers
	messageIncomingConsumer inbound_port.MessageConsumer
	paymentConsumer         inbound_port.MessageConsumer

	// Publishers
	messagePublisher outbound_port.MessagePublisher
}

func NewMessagingManager(rabbitmqURI string) (*MessagingManager, outbound_port.MessagePublisher, error) {
	// Initialize RabbitMQ client
	client, err := message.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	// Setup chat-specific queues
	if err := setupChatQueues(client); err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("failed to setup chat queues: %w", err)
	}

	// Initialize publishers
	messagePublisher := publisherRabbit.NewChatPublisher(client)

	manager := &MessagingManager{
		client:           client,
		messagePublisher: messagePublisher,
	}

	return manager, messagePublisher, nil
}

func (m *MessagingManager) InitializeConsumers(chatService inbound_port.ChatService) {
	m.messageIncomingConsumer = consumerRabbit.NewMessageIncomingConsumer(chatService, m.client)
	m.paymentConsumer = consumerRabbit.NewPaymentConsumer(chatService, m.client)
	log.Println("Consumers initialized successfully")
}

// StartConsumers starts all message consumers
func (m *MessagingManager) StartConsumers() error {
	log.Println("Starting message consumers...")

	if err := m.messageIncomingConsumer.StartListening(); err != nil {
		return fmt.Errorf("failed to start message incoming consumer: %w", err)
	}

	if err := m.paymentConsumer.StartListening(); err != nil {
		return fmt.Errorf("failed to start payment consumer: %w", err)
	}

	log.Println("All consumers started successfully")
	return nil
}

func (m *MessagingManager) GetMessagePublisher() outbound_port.MessagePublisher {
	return m.messagePublisher
}

// Close gracefully shuts down all messaging components
func (m *MessagingManager) Close() error {
	log.Println("Shutting down messaging manager...")
	m.client.Close()
	return nil
}

func setupChatQueues(rmq *message.RabbitMQ) error {
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
