package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/inbound/http"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/inbound/message"
	dbout "github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/publisher"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/repository"
	"github.com/wnmay/horo/services/payment-service/internal/app"
	"github.com/wnmay/horo/services/payment-service/internal/ports/outbound"
	"github.com/wnmay/horo/shared/config"
	"github.com/wnmay/horo/shared/db"
	"github.com/wnmay/horo/shared/env"
	sharedMessage "github.com/wnmay/horo/shared/message"
)

func main() {
	_ = config.LoadEnv("payment-service")
	port := env.GetString("REST_PORT", "3002") // Different port from order service

	// Initialize database
	gormDB := db.MustOpen()
	
	// Initialize repositories
	personRepo := dbout.NewGormPersonRepository(gormDB)
	paymentRepo := repository.NewRepository(gormDB)
	
	// Run migrations
	if err := paymentRepo.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate payment database:", err)
	}
	
	// Initialize event publisher
	eventPublisher := publisher.NewEventPublisher(rabbit)

	// Initialize services
	personSvc := app.NewService(personRepo)
	paymentSvc := app.NewPaymentService(personRepo, outbound.PaymentRepository(paymentRepo), outbound.PaymentEventPublisher(eventPublisher))

	// Initialize RabbitMQ
	rabbitURL := env.GetString("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	rabbit, err := sharedMessage.NewRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer rabbit.Close()

	// Initialize and start message consumer
	consumer := message.NewConsumer(paymentSvc, rabbit)
	go func() {
		log.Println("Starting message consumer...")
		if err := consumer.StartListening(); err != nil {
			log.Fatal("Failed to start consumer:", err)
		}
	}()

	// Initialize HTTP server
	appFiber := fiber.New()
	http.NewHandler(personSvc, paymentSvc).Register(appFiber)

	// Health check
	appFiber.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "payment-service",
		})
	})

	go func() {
		log.Printf("Payment service REST listening on port :%s", port)
		if err := appFiber.Listen(":" + port); err != nil {
			log.Println("server stopped:", err)
		}
	}()
	waitForSignal()
	_ = appFiber.Shutdown()
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
