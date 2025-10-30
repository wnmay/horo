package main

import (
	"context"
	"log"
	"time"

	"github.com/wnmay/horo/services/api-gateway/internal/messaging/publishers"
	"github.com/wnmay/horo/shared/message"
)

func main() {
	log.Println("Starting API Gateway Message Publisher Test...")

	// RabbitMQ connection string
	rabbitmqURI := "replace-rabbit-url-here"

	// Initialize RabbitMQ client
	rmq, err := message.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer rmq.Close()

	log.Println("Connected to RabbitMQ successfully")

	// Create publisher
	publisher := publishers.NewChatMessagePublisher(rmq)

	// Create context
	ctx := context.Background()

	// Test case 1: Publish a text message
	log.Println("\n=== Test 1: Publishing text message ===")
	err = publisher.PublishMessageIncoming(
		ctx,
		"room-123",
		"user-456",
		"Hello from API Gateway test!",
		"text",
	)
	if err != nil {
		log.Fatalf("Failed to publish text message: %v", err)
	}
	log.Println("✓ Text message published successfully")

	// Wait a bit before next message
	time.Sleep(2 * time.Second)

	// Test case 2: Publish a notification message
	log.Println("\n=== Test 2: Publishing notification message ===")
	err = publisher.PublishMessageIncoming(
		ctx,
		"room-789",
		"user-101",
		"User has joined the chat",
		"notification",
	)
	if err != nil {
		log.Fatalf("Failed to publish notification message: %v", err)
	}
	log.Println("✓ Notification message published successfully")

	// Wait a bit before next message
	time.Sleep(2 * time.Second)

	// Test case 3: Publish multiple messages rapidly
	log.Println("\n=== Test 3: Publishing multiple messages rapidly ===")
	for i := 1; i <= 5; i++ {
		err = publisher.PublishMessageIncoming(
			ctx,
			"room-rapid-test",
			"user-rapid",
			"Rapid test message #"+string(rune(i+'0')),
			"text",
		)
		if err != nil {
			log.Printf("Failed to publish message #%d: %v", i, err)
			continue
		}
		log.Printf("✓ Message #%d published", i)
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("\n=== All tests completed ===")
	log.Println("Check the API Gateway logs to see if messages are received by the consumer")
	log.Println("The consumer should show debug logs like:")
	log.Println("[DEBUG] Received chat message: ...")
}
