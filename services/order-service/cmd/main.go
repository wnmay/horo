package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	
	"github.com/wnmay/horo/services/order-service/internal/adapters/inbound/http"
	inboundMessage "github.com/wnmay/horo/services/order-service/internal/adapters/inbound/message"
	"github.com/wnmay/horo/services/order-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/order-service/internal/adapters/outbound/grpc"
	"github.com/wnmay/horo/services/order-service/internal/adapters/outbound/message"
	"github.com/wnmay/horo/services/order-service/internal/app"
	"github.com/wnmay/horo/services/order-service/internal/ports/outbound"
	sharedDB "github.com/wnmay/horo/shared/db"
	"github.com/wnmay/horo/shared/env"
	sharedMessage "github.com/wnmay/horo/shared/message"
)

func main() {
	// Load environment variables
	if err := env.LoadEnv("order-service"); err != nil {
		log.Fatal("Failed to load env:", err)
	}
	
	port := env.GetString("REST_PORT", "3002")

	// Initialize database
	gormDB := sharedDB.MustOpen()
	
	// Initialize repository
	repo := db.NewRepository(gormDB)
	
	// Run migrations
	if err := repo.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	
	// Use the repo as orderRepo (cast to interface)
	orderRepo := outbound.OrderRepository(repo)

	// Initialize message broker (RabbitMQ)
	rabbitURL := env.GetString("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	log.Printf("Using RabbitMQ URL: %s", rabbitURL)
	rabbit, err := sharedMessage.NewRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer rabbit.Close()

	// Initialize course gRPC client
	courseServiceAddr := env.GetString("COURSE_SERVICE_ADDR", "localhost:50052")
	courseClient, err := grpc.NewCourseClient(courseServiceAddr)
	if err != nil {
		log.Fatal("Failed to initialize course client:", err)
	}
	defer courseClient.Close()

	// Initialize adapters
	eventPublisher := message.NewPublisher(rabbit, courseClient)
	paymentService := grpc.NewPaymentClient()
	
	// Initialize application service
	orderService := app.NewOrderService(orderRepo, eventPublisher, paymentService)
	
	// Initialize HTTP handler
	httpHandler := http.NewHandler(orderService)

	// Initialize and start message consumer for payment success events
	consumer := inboundMessage.NewConsumer(orderService, rabbit)
	go func() {
		log.Println("Starting payment success consumer...")
		if err := consumer.StartListening(); err != nil {
			log.Printf("Failed to start payment success consumer: %v", err)
		}
	}()

	// Initialize Fiber app
	appFiber := fiber.New(fiber.Config{
		AppName: "Order Service",
	})
	
	// Add middleware
	appFiber.Use(logger.New())
	appFiber.Use(cors.New())
	
	// Health check
	appFiber.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "order-service",
		})
	})
	
	// Register routes
	httpHandler.Register(appFiber)

	// Start server
	go func() {
		log.Printf("Order service listening on port :%s", port)
		if err := appFiber.Listen(":" + port); err != nil {
			log.Println("Server stopped:", err)
		}
	}()
	
	// Wait for shutdown signal
	waitForSignal()
	
	// Graceful shutdown
	log.Println("Shutting down order service...")
	if err := appFiber.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
