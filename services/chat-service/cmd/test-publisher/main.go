package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/env"
	"github.com/wnmay/horo/shared/message"
)

type MessageIncomingData struct {
	RoomID   string `json:"roomId"`
	SenderID string `json:"senderId"`
	Content  string `json:"content"`
	Type     string `json:"type"` // text | notification
}

func main() {
	err := env.LoadEnv("chat-service")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	rabbitURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
	// Initialize RabbitMQ connection
	rmq, err := message.NewRabbitMQ(rabbitURI)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rmq.Close()

	log.Println("RabbitMQ connection established successfully")

	// Create test message data
	testMessage := MessageIncomingData{
		RoomID:   "69041cc0f18ba67b3f92717a",
		SenderID: "user-456",
		Content:  "Hello from test publisher!",
		Type:     "text",
	}

	// Marshal the test message data
	messageData, err := json.Marshal(testMessage)
	if err != nil {
		log.Fatalf("Failed to marshal test message: %v", err)
	}

	// Create AMQP message wrapper
	amqpMessage := contract.AmqpMessage{
		OwnerID: testMessage.SenderID,
		Data:    messageData,
	}

	// Publish the message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = rmq.PublishMessage(
		ctx,
		contract.ChatMessageIncomingEvent,
		amqpMessage,
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
	log.Printf("Message details: RoomID=%s, SenderID=%s, Content=%s",
		testMessage.RoomID, testMessage.SenderID, testMessage.Content)

	// Wait a bit to ensure message is published
	time.Sleep(1 * time.Second)
	log.Println("Test message published successfully!")
}
