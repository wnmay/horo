package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	chatRabbit "github.com/wnmay/horo/services/chat-service/internal/adapters/inbound/rabbitmq"
	repository "github.com/wnmay/horo/services/chat-service/internal/adapters/outbound/db"
	service "github.com/wnmay/horo/services/chat-service/internal/app"
	"github.com/wnmay/horo/services/chat-service/internal/config"
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

	// Setup chat service specific queues
	if err := chatRabbit.SetupChatQueues(rmq); err != nil {
		log.Fatalf("Failed to setup chat queues: %v", err)
	}

	mongoURI := config.MongoConfig.MongoCommonConfig.URI
	mongoDBName := config.MongoConfig.MongoCommonConfig.Database
	roomCollectionName := config.MongoConfig.RoomCollectionName
	messageCollectionName := config.MongoConfig.MessageCollectionName

	// Initialize MongoDB repositories
	messageRepo, err := repository.NewMongoMessageRepository(mongoURI, mongoDBName, messageCollectionName)
	if err != nil {
		log.Fatalf("Failed to initialize message repository: %v", err)
	}
	roomRepo, err := repository.NewMongoRoomRepository(mongoURI, mongoDBName, roomCollectionName)
	if err != nil {
		log.Fatalf("Failed to initialize room repository: %v", err)
	}

	// Initialize chat service
	chatService := service.NewChatService(messageRepo, roomRepo)

	// Initialize consumer
	messageIncomingConsumer := chatRabbit.NewMessageIncomingConsumer(chatService, rmq)
	log.Println("Msg Incoming Consumer initialized successfully")

	// Start listening for messages
	go messageIncomingConsumer.StartListening()
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down chat service...")
}
