// internal/routes/router.go
package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/services/api-gateway/internal/clients"
	http_handler "github.com/wnmay/horo/services/api-gateway/internal/handlers/http"
	"github.com/wnmay/horo/services/api-gateway/internal/middleware"
)

type Router struct {
	app         *fiber.App
	grpcClients *clients.GrpcClients
}

func NewRouter(app *fiber.App, grpcClients *clients.GrpcClients) *Router {
	return &Router{
		app:         app,
		grpcClients: grpcClients,
	}
}

func (r *Router) SetupRoutes() {
	// Add global logging middleware
	r.app.Use(func(c *fiber.Ctx) error {
		log.Printf("[API-GW] Incoming request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	r.app.Use(middleware.ResponseWrapper())
	// Health check
	r.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	// API v1 group
	api := r.app.Group("/api")

	// Setup all service routes
	r.setupUserRoutes(api)
	r.setupOrderRoutes(api)
	r.setupPaymentRoutes(api)
	r.setupChatRoutes(api)
	r.setupTestRouter(api)
}

func (r *Router) setupUserRoutes(api fiber.Router) {
	userHandler := http_handler.NewUserHandler()

	users := api.Group("/users")
	users.Post("/register", userHandler.Register)
}

func (r *Router) setupOrderRoutes(api fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware(r.grpcClients)
	orderHandler := http_handler.NewOrderHandler()
	orders := api.Group("/orders")

	orders.Post("/", authMiddleware.AddClaims, orderHandler.CreateOrder)
	orders.Get("/:id", authMiddleware.AddClaims, orderHandler.GetOrder)
	// TO DO: don't use cust id here, use from claims
	orders.Get("/customer/:customerID", authMiddleware.AddClaims, orderHandler.GetOrdersByCustomer)
	orders.Patch("/:id/status", authMiddleware.AddClaims, orderHandler.UpdateOrderStatus)
	orders.Patch("/customer/:id", authMiddleware.AddClaims, orderHandler.MarkCustomerCompleted)
	orders.Patch("/prophet/:id", authMiddleware.AddClaims, orderHandler.MarkProphetCompleted)
}

func (r *Router) setupPaymentRoutes(api fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware(r.grpcClients)
	paymentHandler := http_handler.NewPaymentHandler()

	payments := api.Group("/payments")
	payments.Get("/:id", authMiddleware.AddClaims, paymentHandler.GetPayment)
	payments.Get("/order/:orderID", authMiddleware.AddClaims, paymentHandler.GetPaymentByOrder)
	payments.Put("/:id/complete", authMiddleware.AddClaims, paymentHandler.CompletePayment)
}

func (r *Router) setupChatRoutes(api fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware(r.grpcClients)
	chatHandler := http_handler.NewChatHandler()
	chats := api.Group("/chat")
	chats.Get("/:roomID/messages", authMiddleware.AddClaims, chatHandler.GetMessagesByRoomID)
	chats.Post("/rooms", authMiddleware.AddClaims, chatHandler.CreateRoom)
	chats.Get("/customer/rooms", authMiddleware.AddClaims, chatHandler.GetChatRoomsByCustomerID)
	chats.Get("/prophet/rooms", authMiddleware.AddClaims, chatHandler.GetChatRoomsByProphetID)
}

func (r *Router) setupTestRouter(api fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware(r.grpcClients)
	api.Post("/test-auth", authMiddleware.AddClaims, func(c *fiber.Ctx) error {
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to parse body",
			})
		}
		return c.JSON(fiber.Map{
			"message": "middleware worked!",
			"body":    body,
		})
	})

}
