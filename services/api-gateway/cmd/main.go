package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/wnmay/horo/services/api-gateway/internal/config"
	grpcinfra "github.com/wnmay/horo/services/api-gateway/internal/grpc"
	gw_router "github.com/wnmay/horo/services/api-gateway/internal/router"
)

type APIGateway struct {
	app         *fiber.App
	grpcClients *grpcinfra.GrpcClients
	router      *gw_router.Router
	port        string
}

func NewAPIGateway(cfg *config.Config) (*APIGateway, error) {
	// Initialize gRPC clients with all service addresses
	grpcClients, err := grpcinfra.NewGrpcClients(
		cfg.UserManagementAddr,
	)
	if err != nil {
		return nil, err
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(cors.New())

	// Initialize router
	router := gw_router.NewRouter(app, grpcClients)

	return &APIGateway{
		app:         app,
		grpcClients: grpcClients,
		router:      router,
		port:        cfg.Port,
	}, nil
}

func (gw *APIGateway) Start() error {
	// Setup all routes
	gw.router.SetupRoutes()

	log.Printf("Starting API Gateway on port %s", gw.port)
	return gw.app.Listen(":" + gw.port)
}

func (gw *APIGateway) Shutdown() error {
	log.Println("Shutting down API Gateway...")

	gw.grpcClients.Close()

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
