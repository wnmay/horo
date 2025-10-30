package infrastructure

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	http_handler "github.com/wnmay/horo/services/chat-service/internal/adapters/inbound/http"
)

// SetupFiberApp initializes and configures the Fiber application
func SetupFiberApp(messageHandler *http_handler.MessageHandler) *fiber.App {

	app := fiber.New(fiber.Config{
		AppName:      "Chat Service",
		ServerHeader: "Fiber",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "chat-service",
		})
	})

	// Setup HTTP routes
	SetupHTTPRoutes(app, messageHandler)

	return app
}

// StartFiberServer starts the Fiber server on the specified port
func StartFiberServer(app *fiber.App, port string) error {
	log.Printf("Starting Fiber HTTP server on port %s...", port)
	return app.Listen(fmt.Sprintf(":%s", port))
}
