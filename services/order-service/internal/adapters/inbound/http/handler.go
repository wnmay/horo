package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/services/order-service/internal/ports/inbound"
)

type Handler struct{ svc inbound.OrderService }

func NewHandler(s inbound.OrderService) *Handler { return &Handler{svc: s} }

func (h *Handler) Register(app *fiber.App) {
	app.Post("/person", func(c *fiber.Ctx) error {
		var body struct{ Name string `json:"name"` }
		if err := c.BodyParser(&body); err != nil || body.Name == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name is required")
		}
		p, err := h.svc.Create(body.Name)
		if err != nil { return fiber.NewError(fiber.StatusInternalServerError, err.Error()) }
		return c.Status(fiber.StatusCreated).JSON(p)
	})

	app.Get("/person", func(c *fiber.Ctx) error {
		list, err := h.svc.GetAll()
		if err != nil { return fiber.NewError(fiber.StatusInternalServerError, err.Error()) }
		return c.JSON(list)
	})
}
