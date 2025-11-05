package main

import (
	"context"
	"encoding/json"
	"flag"
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
	// Define command-line flags
	eventType := flag.String("event", "chat", "Event type to publish: chat, order-completed, order-payment-bound, order-paid")
	flag.Parse()

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Route to appropriate test based on event type
	switch *eventType {
	case "chat":
		publishChatMessage(ctx, rmq)
	case "order-completed":
		publishOrderCompletedNotification(ctx, rmq)
	case "order-payment-bound":
		publishOrderPaymentBoundNotification(ctx, rmq)
	case "order-paid":
		publishOrderPaidNotification(ctx, rmq)
	default:
		log.Fatalf("Unknown event type: %s. Valid options: chat, order-completed, order-payment-bound, order-paid", *eventType)
	}

	// Wait a bit to ensure message is published
	time.Sleep(1 * time.Second)
	log.Println("Test message published successfully!")
}

func publishChatMessage(ctx context.Context, rmq *message.RabbitMQ) {
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
	err = rmq.PublishMessage(
		ctx,
		contract.ChatMessageIncomingEvent,
		amqpMessage,
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
	log.Printf("✓ Published chat message - RoomID=%s, SenderID=%s, Content=%s",
		testMessage.RoomID, testMessage.SenderID, testMessage.Content)
}

func publishOrderCompletedNotification(ctx context.Context, rmq *message.RabbitMQ) {
	// Create test order completed data
	testData := message.OrderCompletedData{
		OrderID:     "order-test-123",
		PaymentID:   "payment-test-456",
		OrderStatus: "completed",
		CourseID:    "course-test-789",
		CourseName:  "Test Course: Introduction to Testing",
		RoomID:      "69041cc0f18ba67b3f92717a",
	}

	// Marshal the test data
	messageData, err := json.Marshal(testData)
	if err != nil {
		log.Fatalf("Failed to marshal order completed data: %v", err)
	}

	// Create AMQP message wrapper
	amqpMessage := contract.AmqpMessage{
		OwnerID: "system",
		Data:    messageData,
	}

	// Publish the message
	err = rmq.PublishMessage(
		ctx,
		contract.OrderCompletedEvent,
		amqpMessage,
	)
	if err != nil {
		log.Fatalf("Failed to publish order completed notification: %v", err)
	}
	log.Printf("✓ Published OrderCompleted notification - OrderID=%s, CourseID=%s, CourseName=%s, RoomID=%s",
		testData.OrderID, testData.CourseID, testData.CourseName, testData.RoomID)
}

func publishOrderPaymentBoundNotification(ctx context.Context, rmq *message.RabbitMQ) {
	// Create test order payment bound data
	testData := message.OrderPaymentBoundData{
		OrderID:       "order-test-123",
		PaymentID:     "payment-test-456",
		RoomID:        "69041cc0f18ba67b3f92717a",
		CustomerID:    "customer-test-789",
		OrderStatus:   "pending_payment",
		CourseID:      "course-test-789",
		CourseName:    "Test Course: Payment Processing",
		Amount:        999.99,
		PaymentStatus: "pending",
	}

	// Marshal the test data
	messageData, err := json.Marshal(testData)
	if err != nil {
		log.Fatalf("Failed to marshal order payment bound data: %v", err)
	}

	// Create AMQP message wrapper
	amqpMessage := contract.AmqpMessage{
		OwnerID: "system",
		Data:    messageData,
	}

	// Publish the message
	err = rmq.PublishMessage(
		ctx,
		contract.OrderPaymentBoundEvent,
		amqpMessage,
	)
	if err != nil {
		log.Fatalf("Failed to publish order payment bound notification: %v", err)
	}
	log.Printf("✓ Published OrderPaymentBound notification - OrderID=%s, PaymentID=%s, Amount=%.2f, RoomID=%s",
		testData.OrderID, testData.PaymentID, testData.Amount, testData.RoomID)
}

func publishOrderPaidNotification(ctx context.Context, rmq *message.RabbitMQ) {
	// Create test order paid data
	testData := message.OrderPaidData{
		OrderID:       "order-test-123",
		PaymentID:     "payment-test-456",
		RoomID:        "69041cc0f18ba67b3f92717a",
		CustomerID:    "customer-test-789",
		CourseID:      "course-test-789",
		OrderStatus:   "paid",
		CourseName:    "Test Course: Payment Success",
		Amount:        999.99,
		PaymentStatus: "completed",
	}

	// Marshal the test data
	messageData, err := json.Marshal(testData)
	if err != nil {
		log.Fatalf("Failed to marshal order paid data: %v", err)
	}

	// Create AMQP message wrapper
	amqpMessage := contract.AmqpMessage{
		OwnerID: "system",
		Data:    messageData,
	}

	// Publish the message
	err = rmq.PublishMessage(
		ctx,
		contract.OrderPaidEvent,
		amqpMessage,
	)
	if err != nil {
		log.Fatalf("Failed to publish order paid notification: %v", err)
	}
	log.Printf("✓ Published OrderPaid notification - OrderID=%s, PaymentID=%s, Amount=%.2f, PaymentStatus=%s, RoomID=%s",
		testData.OrderID, testData.PaymentID, testData.Amount, testData.PaymentStatus, testData.RoomID)
}
