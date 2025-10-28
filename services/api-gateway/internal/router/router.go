// internal/routes/router.go
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/services/api-gateway/internal/client"
	http_handler "github.com/wnmay/horo/services/api-gateway/internal/http_handler/http"
	"github.com/wnmay/horo/services/api-gateway/internal/middleware"
)

type Router struct {
	app         *fiber.App
	grpcClients *client.GrpcClients
}

func NewRouter(app *fiber.App, grpcClients *client.GrpcClients) *Router {
	return &Router{
		app:         app,
		grpcClients: grpcClients,
	}
}

func (r *Router) SetupRoutes() {
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
	orders.Put("/:id/status", authMiddleware.AddClaims, orderHandler.UpdateOrderStatus)
}

func (r *Router) setupPaymentRoutes(api fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware(r.grpcClients)
	paymentHandler := http_handler.NewPaymentHandler()

	payments := api.Group("/payments")
	payments.Get("/:id", authMiddleware.AddClaims, paymentHandler.GetPayment)
	payments.Get("/order/:orderID", authMiddleware.AddClaims, paymentHandler.GetPaymentByOrder)
	payments.Put("/:id/complete", authMiddleware.AddClaims, paymentHandler.CompletePayment)
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
