package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	http_handler "github.com/wnmay/horo/services/chat-service/internal/adapters/inbound/http"
	repository "github.com/wnmay/horo/services/chat-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/chat-service/internal/app"
	"github.com/wnmay/horo/services/chat-service/internal/config"
	"github.com/wnmay/horo/services/chat-service/internal/infrastructure"
	"github.com/wnmay/horo/services/chat-service/internal/messaging"
	"github.com/wnmay/horo/shared/env"
)

func main() {
	err := env.LoadEnv("chat-service")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	config := config.LoadConfig()

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

	// Initialize messaging manager (handles RabbitMQ client, queues, consumers, and publishers)
	messagingManager, messagePublisher, err := messaging.NewMessagingManager(config.RabbitURI)
	if err != nil {
		log.Fatalf("Failed to initialize messaging manager: %v", err)
	}
	defer messagingManager.Close()

	// Create chat service with repositories and publisher
	chatService := service.NewChatService(messageRepo, roomRepo, messagePublisher)

	// Initialize consumers now that chat service is ready
	messagingManager.InitializeConsumers(chatService)

	// Initialize Fiber HTTP server
	messageHandler := http_handler.NewMessageHandler(chatService)
	app := infrastructure.SetupFiberApp(messageHandler)

	// Start Fiber HTTP server in a goroutine
	go func() {
		if err := infrastructure.StartFiberServer(app, config.HTTPPort); err != nil {
			log.Printf("Fiber server error: %v", err)
		}
	}()

	// Start all message consumers
	if err := messagingManager.StartConsumers(); err != nil {
		log.Fatalf("Failed to start consumers: %v", err)
	}

	log.Println("All consumers are listening and HTTP server is running...")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down chat service...")

	// Gracefully shutdown Fiber server
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber server: %v", err)
	}
}
