package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/shared/config"
	"github.com/wnmay/horo/shared/env"
	sharedDB "github.com/wnmay/horo/shared/db"
)

func main() {
	_ = config.LoadEnv("payment-service")
	port := env.GetString("REST_PORT", "3002")

	log.Println("Starting payment service...")

	// Initialize database
	gormDB := sharedDB.MustOpen()
	log.Println("Database connected successfully")

	// Initialize payment repository (this will auto-migrate the table)
	paymentRepo := db.NewGormPaymentRepository(gormDB)
	log.Printf("Payment table migrated successfully, repository: %v", paymentRepo != nil)

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