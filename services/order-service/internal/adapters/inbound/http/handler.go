package http

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain"
	"github.com/wnmay/horo/services/order-service/internal/ports/inbound"
)

type Handler struct {
	orderService inbound.OrderService
}

func NewHandler(orderService inbound.OrderService) *Handler {
	return &Handler{
		orderService: orderService,
	}
}

func (h *Handler) Register(app *fiber.App) {
	api := app.Group("/api")
	orders := api.Group("/orders")

	// Public routes
	orders.Get("/", h.GetOrders)
	orders.Get("/:id", h.GetOrderByID)

	// require authentication
	orders.Post("/", h.AuthMiddleware, h.CreateOrder)
	orders.Get("/customer/:customerID", h.AuthMiddleware, h.GetOrdersByCustomer)
	orders.Patch("/:id/status", h.AuthMiddleware, h.UpdateOrderStatus)
	orders.Patch("/customer/:id", h.AuthMiddleware, h.MarkCustomerCompleted)
	orders.Patch("/prophet/:id", h.AuthMiddleware, h.MarkProphetCompleted)
}

// AuthMiddleware validates user identity from headers injected by API Gateway
func (h *Handler) AuthMiddleware(c *fiber.Ctx) error {
	// Read user ID from header injected by API Gateway
	userID := c.Get("X-User-Id")
	
	// If header is missing, try to extract from Bearer token (for direct testing)
	if userID == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			// For testing: use a mock user ID
			// In production, this would verify the token with Firebase
			userID = "test-user-from-token"
			log.Printf("Warning: Using mock user ID for direct API testing. Use API Gateway in production.")
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated - X-User-Uid header missing",
			})
		}
	}

	// Optionally read other user info
	email := c.Get("X-User-Email")
	role := c.Get("X-User-Role")

	// Store user info in context for use in handlers
	c.Locals("userID", userID)
	c.Locals("email", email)
	c.Locals("role", role)
	
	return c.Next()
}

type CreateOrderRequest struct {
	CourseID string `json:"courseId" validate:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.CourseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Course ID is required",
		})
	}

	if req.RoomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Room ID is required",
		})
	}

	// Create command with authenticated user ID
	cmd := inbound.CreateOrderCommand{
		CustomerID: userID, // From JWT token
		CourseID:   req.CourseID,
		RoomID:     req.RoomID,
	}

	// Call service
	order, err := h.orderService.CreateOrder(c.Context(), cmd)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}
func (h *Handler) GetOrders(c *fiber.Ctx) error {
	orders, err := h.orderService.GetOrders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(orders)
}

func (h *Handler) GetOrderByID(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	order, err := h.orderService.GetOrderByID(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(order)
}

func (h *Handler) GetOrdersByCustomer(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get customer ID from params
	customerID := c.Params("customerID")
	if customerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Customer ID is required",
		})
	}

	// Check if user is requesting their own orders
	if customerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only access your own orders",
		})
	}

	orders, err := h.orderService.GetOrdersByCustomer(c.Context(), customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(orders)
}

func (h *Handler) UpdateOrderStatus(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse status
	status := req.Status

	if err := h.orderService.UpdateOrderStatus(c.Context(), orderID, domain.OrderStatus(status)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
	})
}

func (h *Handler) MarkCustomerCompleted(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	// Get order to verify ownership
	order, err := h.orderService.GetOrderByID(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	// Verify user is the customer
	if order.CustomerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only mark your own orders as completed",
		})
	}

	if err := h.orderService.MarkCustomerCompleted(c.Context(), orderID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get updated order to return the new status
	updatedOrder, err := h.orderService.GetOrderByID(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve updated order",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order marked as completed by customer",
		"order":   updatedOrder,
	})
}

func (h *Handler) MarkProphetCompleted(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	// TODO: Add logic to verify user is the prophet for this course
	// For now, just mark as completed

	if err := h.orderService.MarkProphetCompleted(c.Context(), orderID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get updated order to return the new status
	updatedOrder, err := h.orderService.GetOrderByID(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve updated order",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order marked as completed by prophet",
		"order":   updatedOrder,
	})
}
