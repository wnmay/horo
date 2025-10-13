package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("Payment service starting...")
	
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "payment-service",
		})
	})
	
	log.Println("Payment service listening on port 3002...")
	log.Fatal(app.Listen(":3002"))
}