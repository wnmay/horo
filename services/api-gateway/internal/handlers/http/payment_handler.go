// internal/handlers/http/payment_handler.go
package http_handler

import (
	"fmt"
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
	return ProxyRequest(c,  h.client, "GET", h.paymentServiceURL, fmt.Sprintf("/api/payments/%s", id))
}

func (h *PaymentHandler) GetPaymentByOrder(c *fiber.Ctx) error {
	orderID := c.Params("orderID")
	return ProxyRequest(c,  h.client, "GET", h.paymentServiceURL, fmt.Sprintf("/api/payments/order/%s", orderID))
}

func (h *PaymentHandler) CompletePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	return ProxyRequest(c,  h.client, "PUT", h.paymentServiceURL, fmt.Sprintf("/api/payments/%s/complete", id))
}
