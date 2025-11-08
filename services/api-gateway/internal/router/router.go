// internal/routes/router.go
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/wnmay/horo/services/api-gateway/internal/config"
	http_handler "github.com/wnmay/horo/services/api-gateway/internal/handlers/http"
	ws_handler "github.com/wnmay/horo/services/api-gateway/internal/handlers/ws"
	"github.com/wnmay/horo/services/api-gateway/internal/messaging/publishers"
	"github.com/wnmay/horo/services/api-gateway/internal/middleware"
	gwWS "github.com/wnmay/horo/services/api-gateway/internal/websocket"
	"github.com/wnmay/horo/shared/message"
)

type Router struct {
	app            *fiber.App
	rmq            *message.RabbitMQ
	hub            *gwWS.Hub
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(app *fiber.App, cfg *config.Config, rmq *message.RabbitMQ) *Router {
	return &Router{
		app:            app,
		rmq:            rmq,
		hub:            gwWS.NewHub(),
		authMiddleware: middleware.NewAuthMiddleware(cfg.UserManagementServiceURL),
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
	orderHandler := http_handler.NewOrderHandler()
	orders := api.Group("/orders")

	orders.Get("/", r.authMiddleware.AddClaims, orderHandler.GetOrders)
	orders.Post("/", r.authMiddleware.AddClaims, orderHandler.CreateOrder)
	orders.Get("/:id", r.authMiddleware.AddClaims, orderHandler.GetOrderByID)
	orders.Get("/customer/:customerID", r.authMiddleware.AddClaims, orderHandler.GetOrdersByCustomer)
	orders.Patch("/:id/status", r.authMiddleware.AddClaims, orderHandler.UpdateOrderStatus)
	orders.Patch("/customer/:id", r.authMiddleware.AddClaims, orderHandler.MarkCustomerCompleted)
	orders.Patch("/prophet/:id", r.authMiddleware.AddClaims, orderHandler.MarkProphetCompleted)
}

func (r *Router) setupPaymentRoutes(api fiber.Router) {
	paymentHandler := http_handler.NewPaymentHandler()

	payments := api.Group("/payments")
	payments.Get("/:id", r.authMiddleware.AddClaims, paymentHandler.GetPayment)
	payments.Get("/order/:orderID", r.authMiddleware.AddClaims, paymentHandler.GetPaymentByOrder)
	payments.Put("/:id/complete", r.authMiddleware.AddClaims, paymentHandler.CompletePayment)
	payments.Get("/balance", r.authMiddleware.AddClaims, paymentHandler.GetProphetBalance)
}

func (r *Router) setupChatRoutes(api fiber.Router) {
	chatHandler := http_handler.NewChatHandler()
	chats := api.Group("/chat")
	chats.Get("/:roomID/messages", r.authMiddleware.AddClaims, chatHandler.GetMessagesByRoomID)
	chats.Post("/rooms", r.authMiddleware.AddClaims, chatHandler.CreateRoom)
	chats.Get("/customer/rooms", r.authMiddleware.AddClaims, chatHandler.GetChatRoomsByCustomerID)
	chats.Get("/prophet/rooms", r.authMiddleware.AddClaims, chatHandler.GetChatRoomsByProphetID)
	chats.Get("/user/rooms", r.authMiddleware.AddClaims, chatHandler.GetChatRoomsByUserID)
}

func (r *Router) setupCourseRoutes(api fiber.Router) {
	courseHandler := http_handler.NewCourseHandler()
	authMiddleware := r.authMiddleware

	courses := api.Group("/courses")

	courses.Post("/", authMiddleware.AddClaims, courseHandler.CreateCourse)
	courses.Get("/:id", authMiddleware.AddClaims, courseHandler.GetCourseByID)
	courses.Get("/prophet/:prophetId/courses", authMiddleware.AddClaims, courseHandler.ListCoursesByProphet)
	courses.Patch("/:id", authMiddleware.AddClaims, courseHandler.UpdateCourse)
	courses.Patch("/delete/:id", authMiddleware.AddClaims, courseHandler.DeleteCourse)
	courses.Get("/", authMiddleware.AddClaims, courseHandler.FindCoursesByFilter)
	courses.Post("/:courseId/review", authMiddleware.AddClaims, courseHandler.CreateReview)
	courses.Get("/review/:id", authMiddleware.AddClaims, courseHandler.GetReviewByID)
	courses.Get("/:courseId/reviews", authMiddleware.AddClaims, courseHandler.ListReviewsByCourse)
}

func (r *Router) setupTestRouter(api fiber.Router) {
	api.Post("/test-auth", r.authMiddleware.AddClaims, func(c *fiber.Ctx) error {
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
	chatPublisher := publishers.NewChatMessagePublisher(r.rmq)
	chatWsHandler := ws_handler.NewChatWSHandler(r.hub, chatPublisher)

	r.app.Use("/ws/chat", r.authMiddleware.AddClaimsWS, func(c *fiber.Ctx) error {
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
