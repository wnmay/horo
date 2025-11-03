package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/wnmay/horo/services/payment-service/internal/adapters/inbound/http"
	inboundMessage "github.com/wnmay/horo/services/payment-service/internal/adapters/inbound/message"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/message"
	"github.com/wnmay/horo/services/payment-service/internal/app"
	sharedDB "github.com/wnmay/horo/shared/db"
	"github.com/wnmay/horo/shared/env"
	sharedMessage "github.com/wnmay/horo/shared/message"
)

func main() {
	_ = env.LoadEnv("payment-service")
	port := env.GetString("REST_PORT", "3001")

	log.Println("Starting payment service...")

	// Initialize database
	gormDB := sharedDB.MustOpen()
	log.Println("Database connected successfully")

	// Initialize payment repository (this will auto-migrate the table)
	paymentRepo := db.NewGormPaymentRepository(gormDB)

	// Initialize RabbitMQ
	rabbitURL := env.GetString("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	log.Printf("Connecting to RabbitMQ: %s", rabbitURL)
	rabbit, err := sharedMessage.NewRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer rabbit.Close()
	log.Println("RabbitMQ connected successfully")
	
	// Initialize publisher
	eventPublisher := message.NewPublisher(rabbit)
	
	// Initialize application service
	paymentService := app.NewPaymentService(paymentRepo, eventPublisher)
	
	// Initialize consumer
	consumer := inboundMessage.NewConsumer(paymentService, rabbit)
	go func() {
		log.Println("Starting order created consumer...")
		if err := consumer.StartListening(); err != nil {
			log.Printf("Failed to start order created consumer: %v", err)
		}
	}()
	
	// Initialize HTTP server
	httpHandler := http.NewHandler(paymentService)
	
	// Initialize fiber app
	appFiber := fiber.New(fiber.Config{
		AppName: "Payment Service",
	})
	
	// Add middleware
	appFiber.Use(cors.New())
	
	// Health check
	appFiber.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "payment-service",
		})
	})

	// Register routes
	httpHandler.Register(appFiber)

	// Start server
	go func() {
		log.Printf("Payment service listening on port :%s", port)
		if err := appFiber.Listen(":" + port); err != nil {
			log.Println("Server stopped:", err)
		}
	}()
	
	// Wait for shutdown signal
	waitForSignal()
	
	// Graceful shutdown
	log.Println("Shutting down payment service...")
	if err := appFiber.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}