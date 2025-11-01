// internal/handlers/http/order_handler.go
package http_handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/shared/env"
)

type OrderHandler struct {
	orderServiceURL string
	client          *http.Client
}

func NewOrderHandler() *OrderHandler {
	orderServiceURL := env.GetString("ORDER_SERVICE_URL", "localhost:3002")

	return &OrderHandler{
		orderServiceURL: orderServiceURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "POST", h.orderServiceURL, "/api/orders/")
}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	return ProxyRequest(c, h.client, "GET", h.orderServiceURL, fmt.Sprintf("/api/orders/%s", id))
}

// TO DO: change cust id to user id
func (h *OrderHandler) GetOrdersByCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customerID")
	return ProxyRequest(c, h.client, "GET", h.orderServiceURL, fmt.Sprintf("/api/orders/customer/%s", customerID))
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	return ProxyRequest(c, h.client, "PUT", h.orderServiceURL, fmt.Sprintf("/api/orders/%s/status", id))
}
