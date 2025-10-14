package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/shared/config"
	"github.com/wnmay/horo/shared/env"
	sharedDB "github.com/wnmay/horo/shared/db"
	dbout "github.com/wnmay/horo/services/payment-service/internal/adapters/outbound/db"
)

func main() {
	_ = config.LoadEnv("payment-service")
	port := env.GetString("REST_PORT", "3002")

	// Initialize database
	gormDB := sharedDB.MustOpen()
	log.Println("Database connected:", gormDB != nil)

	// Initialize payment repository and auto-migrate
	paymentRepo := dbout.NewGormPaymentRepository(gormDB)
	log.Println("Payment repository initialized and migrated")

	// Initialize HTTP server
	app := fiber.New()

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "payment-service",
		})
	})

	go func() {
		log.Printf("Payment service REST listening on port :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Println("server stopped:", err)
		}
	}()
	waitForSignal()
	_ = app.Shutdown()
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}