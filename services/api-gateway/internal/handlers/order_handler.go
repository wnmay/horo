// internal/handlers/order_handler.go
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	return h.proxyRequest(c, "POST", "/api/orders/")
}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	return h.proxyRequest(c, "GET", fmt.Sprintf("/api/orders/%s", id))
}

// TO DO: change cust id to user id
func (h *OrderHandler) GetOrdersByCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customerID")
	return h.proxyRequest(c, "GET", fmt.Sprintf("/api/orders/customer/%s", customerID))
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	return h.proxyRequest(c, "PUT", fmt.Sprintf("/api/orders/%s/status", id))
}

func (h *OrderHandler) proxyRequest(c *fiber.Ctx, method, path string) error {
	targetURL := h.orderServiceURL + path

	// Get query parameters
	query := c.Request().URI().QueryString()
	if len(query) > 0 {
		targetURL += "?" + string(query)
	}

	// Create request body
	var body io.Reader
	if method == "POST" || method == "PUT" {
		body = bytes.NewReader(c.Body())
	}

	// Create HTTP request
	req, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create request",
		})
	}

	// Copy headers
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	// Make request
	resp, err := h.client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "failed to reach order service",
		})
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read response",
		})
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	// Return response
	c.Status(resp.StatusCode)

	// Try to parse as JSON, if fails return as raw
	var jsonResp interface{}
	if err := json.Unmarshal(respBody, &jsonResp); err == nil {
		return c.JSON(jsonResp)
	}

	return c.Send(respBody)
}
