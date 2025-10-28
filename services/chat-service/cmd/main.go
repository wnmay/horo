package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	consumerRabbit "github.com/wnmay/horo/services/chat-service/internal/adapters/inbound/rabbitmq"
	repository "github.com/wnmay/horo/services/chat-service/internal/adapters/outbound/db"
	publisherRabbit "github.com/wnmay/horo/services/chat-service/internal/adapters/outbound/rabbitmq"
	service "github.com/wnmay/horo/services/chat-service/internal/app"
	"github.com/wnmay/horo/services/chat-service/internal/config"
	"github.com/wnmay/horo/services/chat-service/internal/infrastructure"
	"github.com/wnmay/horo/shared/env"
	"github.com/wnmay/horo/shared/message"
)

func main() {
	err := env.LoadEnv("chat-service")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	config := config.LoadConfig()

	// Initialize RabbitMQ connection (sets up centralized infrastructure)
	rmq, err := message.NewRabbitMQ(config.RabbitURI)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rmq.Close()

	if err := infrastructure.SetupChatQueues(rmq); err != nil {
		log.Fatalf("Failed to setup chat queues: %v", err)
	}

	mongoURI := config.MongoConfig.MongoCommonConfig.URI
	mongoDBName := config.MongoConfig.MongoCommonConfig.Database
	roomCollectionName := config.MongoConfig.RoomCollectionName
	messageCollectionName := config.MongoConfig.MessageCollectionName

	// Initialize MongoDB (single connection pool for all repositories)
	mongoDB, mongoClient, err := infrastructure.SetupMongoDB(mongoURI, mongoDBName)
	if err != nil {
		log.Fatalf("Failed to setup MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Initialize MongoDB repositories (sharing the same client)
	messageRepo := repository.NewMongoMessageRepository(mongoDB, messageCollectionName)
	roomRepo := repository.NewMongoRoomRepository(mongoDB, roomCollectionName)

	messagePublisher := publisherRabbit.NewChatPublisher(rmq)

	chatService := service.NewChatService(messageRepo, roomRepo, messagePublisher)

	// Initialize consumer
	messageIncomingConsumer := consumerRabbit.NewMessageIncomingConsumer(chatService, rmq)
	paymentConsumer := consumerRabbit.NewPaymentConsumer(chatService, rmq)
	log.Println("Consumers initialized successfully")

	// Start listening for messages from both consumers
	go func() {
		if err := messageIncomingConsumer.StartListening(); err != nil {
			log.Printf("Message incoming consumer error: %v", err)
		}
	}()

	go func() {
		if err := paymentConsumer.StartListening(); err != nil {
			log.Printf("Payment consumer error: %v", err)
		}
	}()

	log.Println("All consumers are listening...")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down chat service...")
}
