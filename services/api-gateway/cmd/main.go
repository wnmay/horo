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
	"github.com/wnmay/horo/services/api-gateway/internal/messaging"
	gw_router "github.com/wnmay/horo/services/api-gateway/internal/router"
	"github.com/wnmay/horo/shared/env"
)

type APIGateway struct {
	app              *fiber.App
	grpcClients      *client.GrpcClients
	messagingManager *messaging.MessagingManager
	router           *gw_router.Router
	port             string
	cfg              *config.Config
}

const (
	service_name = "api-gateway"
)

func NewAPIGateway(cfg *config.Config) (*APIGateway, error) {
	// Initialize gRPC clients with all service addresses
	grpcClients, err := client.NewGrpcClients(
		cfg.UserManagementAddr,
		cfg.ChatAddr,
	)
	if err != nil {
		return nil, err
	}

	// Initialize messaging manager (handles RabbitMQ client, consumers, and publishers)
	messagingManager, err := messaging.NewMessagingManager(cfg.RabbitMQURI)
	log.Println("RabbitMQ URI:", cfg.RabbitMQURI)
	if err != nil {
		grpcClients.Close()
		return nil, err
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(cors.New())

	// Initialize router
	router := gw_router.NewRouter(app, grpcClients,messagingManager.RabbitMQ())

	return &APIGateway{
		app:              app,
		grpcClients:      grpcClients,
		messagingManager: messagingManager,
		router:           router,
		port:             cfg.Port,
		cfg:              cfg,
	}, nil
}

func (gw *APIGateway) Start() error {
	// Setup all routes
	gw.router.SetupRoutes()
	hub := gw.router.GetHub()

	gw.messagingManager.StartChatConsumer(hub)

	log.Printf("Starting API Gateway on port %s", gw.port)
	return gw.app.Listen(":" + gw.port)
}

func (gw *APIGateway) Shutdown() error {
	log.Println("Shutting down API Gateway...")

	gw.grpcClients.Close()
	gw.messagingManager.Close()

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
