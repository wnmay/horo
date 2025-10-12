// internal/routes/router.go
package router

import (
	"github.com/gofiber/fiber/v2"
	grpcinfra "github.com/wnmay/horo/services/api-gateway/internal/grpc"
	"github.com/wnmay/horo/services/api-gateway/internal/handlers"
)

type Router struct {
	app         *fiber.App
	grpcClients *grpcinfra.GrpcClients
}

func NewRouter(app *fiber.App, grpcClients *grpcinfra.GrpcClients) *Router {
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
}

func (r *Router) setupUserRoutes(api fiber.Router) {
	userHandler := handlers.NewUserHandler(r.grpcClients)

	users := api.Group("/users")
	users.Post("/register", userHandler.Register)
}
