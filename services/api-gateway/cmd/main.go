package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	client "github.com/wnmay/horo/services/api-gateway/internal/clients"
	"github.com/wnmay/horo/services/api-gateway/internal/config"
	"github.com/wnmay/horo/services/api-gateway/internal/messaging/consumers"
	gw_router "github.com/wnmay/horo/services/api-gateway/internal/router"
	"github.com/wnmay/horo/shared/env"
	shared_message "github.com/wnmay/horo/shared/message"
)

type APIGateway struct {
	app            *fiber.App
	grpcClients    *client.GrpcClients
	rabbitmqClient *shared_message.RabbitMQ
	chatConsumer   *consumers.ChatMessageConsumer
	router         *gw_router.Router
	port           string
}

const (
	service_name = "api-gateway"
)

func NewAPIGateway(cfg *config.Config) (*APIGateway, error) {
	// Initialize gRPC clients with all service addresses
	grpcClients, err := client.NewGrpcClients(
		cfg.UserManagementAddr,
	)
	if err != nil {
		return nil, err
	}

	// Initialize RabbitMQ client
	rabbitmqClient, err := shared_message.NewRabbitMQ(cfg.RabbitMQURI)
	if err != nil {
		grpcClients.Close()
		return nil, err
	}

	// Initialize chat message consumer
	chatConsumer := consumers.NewChatMessageConsumer(rabbitmqClient)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(cors.New())

	// Initialize router
	router := gw_router.NewRouter(app, grpcClients)

	return &APIGateway{
		app:            app,
		grpcClients:    grpcClients,
		rabbitmqClient: rabbitmqClient,
		chatConsumer:   chatConsumer,
		router:         router,
		port:           cfg.Port,
	}, nil
}

func (gw *APIGateway) Start() error {
	// Setup all routes
	gw.router.SetupRoutes()

	// Start consuming chat messages
	if err := gw.chatConsumer.StartListening(); err != nil {
		return err
	}

	log.Printf("Starting API Gateway on port %s", gw.port)
	return gw.app.Listen(":" + gw.port)
}

func (gw *APIGateway) Shutdown() error {
	log.Println("Shutting down API Gateway...")

	gw.grpcClients.Close()
	gw.rabbitmqClient.Close()

	return gw.app.Shutdown()
}
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}

func main() {
	_ = env.LoadEnv(service_name)
	cfg := config.LoadConfig()
	// Create API Gateway
	gateway, err := NewAPIGateway(cfg)
	if err != nil {
		log.Fatalf("Failed to create API Gateway: %v", err)
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received shutdown signal")
		if err := gateway.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		os.Exit(0)
	}()

	// Start server
	if err := gateway.Start(); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
