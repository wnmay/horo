package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wnmay/horo/services/order-service/internal/domain/entity"
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
	api := app.Group("/api/v1")
	orders := api.Group("/orders")

	orders.Post("/", h.CreateOrder)
	orders.Get("/:id", h.GetOrder)
	orders.Get("/customer/:customerID", h.GetOrdersByCustomer)
	orders.Put("/:id/status", h.UpdateOrderStatus)
}

type Claims struct {
	CustomerID string `json:"customerId" validate:"required"`
}

type CreateOrderRequest struct {
	Claims   Claims `json:"claims" validate:"required"`
	CourseID string `json:"course_id" validate:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Debug logging to see what we received
	fmt.Printf("Received request: %+v\n", req)
	fmt.Printf("Claims: %+v\n", req.Claims)
	fmt.Printf("CustomerID: '%s'\n", req.Claims.CustomerID)
	fmt.Printf("CourseID: '%s'\n", req.CourseID)

	// Validate required fields
	if req.Claims.CustomerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Customer ID is required in claims",
		})
	}

	if req.CourseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Course ID is required",
		})
	}

	// Parse course UUID
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid course ID format",
		})
	}

	// Create command
	cmd := inbound.CreateOrderCommand{
		CustomerID: req.Claims.CustomerID, // Firebase userId as string
		CourseID:   courseID,
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

func (h *Handler) GetOrder(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	order, err := h.orderService.GetOrder(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(order)
}

func (h *Handler) GetOrdersByCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customerID")
	if customerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Customer ID is required",
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

	if err := h.orderService.UpdateOrderStatus(c.Context(), orderID, entity.OrderStatus(status)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
	})
}