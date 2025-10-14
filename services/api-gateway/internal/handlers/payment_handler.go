// internal/handlers/payment_handler.go
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

type PaymentHandler struct {
	paymentServiceURL string
	client            *http.Client
}

func NewPaymentHandler() *PaymentHandler {
	paymentServiceURL := env.GetString("PAYMENT_SERVICE_URL", "http://localhost:3001")
	return &PaymentHandler{
		paymentServiceURL: paymentServiceURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	id := c.Params("id")
	return h.proxyRequest(c, "GET", fmt.Sprintf("/api/payments/%s", id))
}

func (h *PaymentHandler) GetPaymentByOrder(c *fiber.Ctx) error {
	orderID := c.Params("orderID")
	return h.proxyRequest(c, "GET", fmt.Sprintf("/api/payments/order/%s", orderID))
}

func (h *PaymentHandler) CompletePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	return h.proxyRequest(c, "PUT", fmt.Sprintf("/api/payments/%s/complete", id))
}

func (h *PaymentHandler) proxyRequest(c *fiber.Ctx, method, path string) error {
	// Build target URL
	targetURL := h.paymentServiceURL + path

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

	// Copy headers (including Authorization for auth)
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	// Make request
	resp, err := h.client.Do(req)
	if err != nil {
		errstr := fmt.Sprintf("Failed to reach payment service: %v\n", err) // ADD THIS
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": errstr,
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
