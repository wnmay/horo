package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/services/payment-service/internal/ports/inbound"
)

type Handler struct {
	personSvc  inbound.PersonService
	paymentSvc inbound.PaymentService
}

func NewHandler(personSvc inbound.PersonService, paymentSvc inbound.PaymentService) *Handler {
	return &Handler{
		personSvc:  personSvc,
		paymentSvc: paymentSvc,
	}
}

func (h *Handler) Register(app *fiber.App) {
	api := app.Group("/api/v1")
	payments := api.Group("/payments")

	payments.Get("/:id", h.GetPayment)
	payments.Get("/order/:orderID", h.GetPaymentByOrder)
	payments.Put("/:id/complete", h.CompletePayment)
}

func (h *Handler) GetPayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")
	if paymentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment ID is required",
		})
	}

	payment, err := h.paymentSvc.GetPayment(c.Context(), paymentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(payment)
}

func (h *Handler) GetPaymentByOrder(c *fiber.Ctx) error {
	orderID := c.Params("orderID")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	payment, err := h.paymentSvc.GetPaymentByOrderID(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(payment)
}

func (h *Handler) CompletePayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")
	if paymentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment ID is required",
		})
	}

	if err := h.paymentSvc.CompletePayment(c.Context(), paymentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment completed successfully",
		"payment_id": paymentID,
	})
}