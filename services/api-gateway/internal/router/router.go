// internal/routes/router.go
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/wnmay/horo/services/api-gateway/internal/clients"
	http_handler "github.com/wnmay/horo/services/api-gateway/internal/handlers/http"
	ws_handler "github.com/wnmay/horo/services/api-gateway/internal/handlers/ws"
	"github.com/wnmay/horo/services/api-gateway/internal/messaging/publishers"
	"github.com/wnmay/horo/services/api-gateway/internal/middleware"
	gwWS "github.com/wnmay/horo/services/api-gateway/internal/websocket"
	"github.com/wnmay/horo/shared/message"
)

type Router struct {
	app         *fiber.App
	grpcClients *clients.GrpcClients
	rmq         *message.RabbitMQ
	hub         *gwWS.Hub
}

func NewRouter(app *fiber.App, grpcClients *clients.GrpcClients, rmq *message.RabbitMQ) *Router {
	return &Router{
		app:         app,
		grpcClients: grpcClients,
		rmq:         rmq,
		hub:         gwWS.NewHub(),
	}
}

func (r *Router) SetupRoutes() {
	r.app.Use(middleware.ResponseWrapper())

	// Health check
	r.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	r.setupWebsocketRoutes()

	// API v1 group
	api := r.app.Group("/api")

	// Setup all service routes
	r.setupUserRoutes(api)
	r.setupOrderRoutes(api)
	r.setupPaymentRoutes(api)
	r.setupChatRoutes(api)
	r.setupCourseRoutes(api)
	r.setupTestRouter(api)
	r.setupChatRoutes(api)
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

	orders.Get("/", authMiddleware.AddClaims, orderHandler.GetOrders)
	orders.Post("/", authMiddleware.AddClaims, orderHandler.CreateOrder)
	orders.Get("/:id", authMiddleware.AddClaims, orderHandler.GetOrderByID)
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
	chats := api.Group("/chats")
	chats.Get("/:roomID/messages", authMiddleware.AddClaims, chatHandler.GetMessagesByRoomID)
	chats.Post("/rooms", authMiddleware.AddClaims, chatHandler.CreateRoom)
	chats.Get("/customer/rooms", authMiddleware.AddClaims, chatHandler.GetChatRoomsByCustomerID)
	chats.Get("/prophet/rooms", authMiddleware.AddClaims, chatHandler.GetChatRoomsByProphetID)
	chats.Get("/user/rooms", authMiddleware.AddClaims, chatHandler.GetChatRoomsByUserID)
}

func (r *Router) setupCourseRoutes(api fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware(r.grpcClients)
	courseHandler := http_handler.NewCourseHandler()

	courses := api.Group("/courses")

	courses.Post("/", authMiddleware.AddClaims, courseHandler.CreateCourse)
	courses.Get("/:id", authMiddleware.AddClaims, courseHandler.GetCourseByID)
	courses.Get("/prophet/:prophet_id", authMiddleware.AddClaims, courseHandler.ListCoursesByProphet)
	courses.Patch("/:id", authMiddleware.AddClaims, courseHandler.UpdateCourse)
	courses.Patch("/delete/:id", authMiddleware.AddClaims, courseHandler.DeleteCourse)
	courses.Get("/", authMiddleware.AddClaims, courseHandler.FindCoursesByFilter)
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

func (r *Router) setupWebsocketRoutes() {
	authMiddleware := middleware.NewAuthMiddleware(r.grpcClients)
	chatPublisher := publishers.NewChatMessagePublisher(r.rmq)
	chatWsHandler := ws_handler.NewChatWSHandler(r.hub, chatPublisher, r.grpcClients.ChatServiceClient)

	r.app.Use("/ws/chat", authMiddleware.AddClaims, func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	chatWsHandler.RegisterRoutes(r.app)
}

func (r *Router) GetHub() *gwWS.Hub {
	return r.hub
}
