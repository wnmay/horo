package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/shared/config"
	"github.com/wnmay/horo/shared/contract"
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

	// API routes with versioning
	api := app.Group("/api/v1")
	
	// Basic payment endpoints for testing
	api.Get("/payments", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Payments endpoint - ready for implementation",
		})
	})

	// Payment completion endpoint
	api.Put("/payments/:id/complete", func(c *fiber.Ctx) error {
		paymentID := c.Params("id")
		if paymentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Payment ID is required",
			})
		}

		log.Printf("Processing payment completion for ID: %s", paymentID)

		// Get the payment from database
		payment, err := paymentRepo.GetByID(c.Context(), paymentID)
		if err != nil {
			log.Printf("Failed to get payment: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Payment not found",
			})
		}

		// Update payment status to completed
		payment.Status = "COMPLETED"
		payment.CreatedAt = time.Now() // Update timestamp

		if err := paymentRepo.Update(c.Context(), payment); err != nil {
			log.Printf("Failed to update payment: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update payment status",
			})
		}

		log.Printf("Payment %s marked as completed", paymentID)

		// Publish payment completed event to notify order service
		go func() {
			ctx := context.Background()
			
			// Create payment completion message
			paymentCompletedData := map[string]interface{}{
				"payment_id": payment.PaymentID,
				"order_id":   payment.OrderID,
				"status":     "COMPLETED",
				"amount":     payment.Amount,
			}

			dataBytes, err := json.Marshal(paymentCompletedData)
			if err != nil {
				log.Printf("Failed to marshal payment completion data: %v", err)
				return
			}

			// Create AMQP message using proper contract structure
			amqpMessage := contract.AmqpMessage{
				OwnerID: payment.OrderID,
				Data:    dataBytes,
			}

			if err := rabbit.PublishMessage(ctx, "payment.completed", amqpMessage); err != nil {
				log.Printf("Failed to publish payment completion event: %v", err)
			} else {
				log.Printf("Published payment completed event for order: %s", payment.OrderID)
			}
		}()

		return c.JSON(fiber.Map{
			"message":    "Payment completed successfully",
			"payment_id": paymentID,
			"order_id":   payment.OrderID,
			"status":     "COMPLETED",
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
			
			// Payment creation handler
			handler := func(ctx context.Context, delivery amqp.Delivery) error {
				log.Printf("Received order event: %s", string(delivery.Body))
				
				// Parse the AMQP message
				var amqpMessage struct {
					OwnerID string `json:"ownerId"`
					Data    string `json:"data"` // Base64 encoded JSON
				}
				
				if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
					log.Printf("Failed to unmarshal AMQP message: %v", err)
					return err
				}
				
				// Decode base64 data
				decodedData, err := base64.StdEncoding.DecodeString(amqpMessage.Data)
				if err != nil {
					log.Printf("Failed to decode base64 data: %v", err)
					return err
				}
				
				log.Printf("Decoded order data: %s", string(decodedData))
				
				// Parse the order data
				var orderData struct {
					OrderID    string  `json:"order_id"`
					CustomerID string  `json:"customer_id"`
					Amount     float64 `json:"amount"`
					Status     string  `json:"status"`
				}
				
				if err := json.Unmarshal(decodedData, &orderData); err != nil {
					log.Printf("Failed to unmarshal order data: %v", err)
					return err
				}
				
				// Default payment amount (in real scenario, this would be fetched based on course pricing)
				defaultAmount := 149.99
				
				log.Printf("Creating payment for order: %s, amount: %.2f", orderData.OrderID, defaultAmount)
				
				// Create payment using our domain entity
				payment := &domain.Payment{
					PaymentID: uuid.New().String(),
					OrderID:   orderData.OrderID,
					Amount:    defaultAmount,
					Status:    "pending",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				
				// Save payment to database
				if err := paymentRepo.Create(ctx, payment); err != nil {
					log.Printf("Failed to create payment: %v", err)
					return err
				}
				
				log.Printf("Payment created successfully: %s for order: %s", payment.PaymentID, orderData.OrderID)
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
