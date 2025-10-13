package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/shared/config"
	"github.com/wnmay/horo/shared/env"
	sharedDB "github.com/wnmay/horo/shared/db"
	sharedMessage "github.com/wnmay/horo/shared/message"
)

func main() {
	_ = config.LoadEnv("payment-service")
	port := env.GetString("REST_PORT", "3001")

	log.Println("Starting payment service...")

	// Initialize database
	gormDB := sharedDB.MustOpen()
	log.Println("Database connected successfully")

	// Initialize payment repository (this will auto-migrate the table)
	paymentRepo := db.NewGormPaymentRepository(gormDB)
	log.Printf("Payment table migrated successfully, repository: %v", paymentRepo != nil)

	// Initialize RabbitMQ
	rabbitURL := env.GetString("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	log.Printf("Connecting to RabbitMQ: %s", rabbitURL)
	rabbit, err := sharedMessage.NewRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer rabbit.Close()
	log.Println("RabbitMQ connected successfully")

	// Declare the queue for receiving order created events
	if err := rabbit.DeclareQueue("create_payment_queue", "order.created"); err != nil {
		log.Fatal("Failed to declare create payment queue:", err)
	}
	log.Println("Payment queue declared successfully")

	// Initialize HTTP server  
	app := fiber.New()

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "payment-service",
		})
	})

	// Basic payment endpoints for testing
	app.Get("/payments", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Payments endpoint - ready for implementation",
		})
	})

	go func() {
		log.Printf("Payment service listening on port :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Println("Server stopped:", err)
		}
	}()
	
	// Start listening for order created events (RabbitMQ already initialized above)
	go func() {
			queueName := "payment_queue"
			routingKey := "order.created"
			
			// Declare queue
			if err := rabbit.DeclareQueue(queueName, routingKey); err != nil {
				log.Printf("Failed to declare queue: %v", err)
				return
			}
			
			log.Printf("Listening for order events on queue: %s", queueName)
			
			// Simple message handler for now
			handler := func(ctx context.Context, delivery amqp.Delivery) error {
				log.Printf("Received order event: %s", string(delivery.Body))
				
				// TODO: Parse the order data and create payment in database
				// For now, just log that we received it
				log.Println("Payment should be created here!")
				
				return nil
			}
			
			if err := rabbit.ConsumeMessages(queueName, handler); err != nil {
				log.Printf("Failed to start consuming messages: %v", err)
			}
		}()

	log.Println("Payment service started successfully! Press Ctrl+C to stop.")
	waitForSignal()
	_ = app.Shutdown()
	log.Println("Payment service stopped.")
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
