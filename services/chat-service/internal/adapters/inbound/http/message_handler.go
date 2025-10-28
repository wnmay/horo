package http_handler

import (
	"github.com/gofiber/fiber/v2"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
)

type MessageHandler struct {
	chatService inbound_port.ChatService
}

type CreateRoomRequest struct {
	CourseID   string `json:"courseID"`
	CustomerID string `json:"customerID"`
}

func NewMessageHandler(chatService inbound_port.ChatService) *MessageHandler {
	return &MessageHandler{
		chatService: chatService,
	}
}

func (h *MessageHandler) GetMessagesByRoomID(c *fiber.Ctx) error {
	roomID := c.Params("roomID")
	messages, err := h.chatService.GetMessagesByRoomID(c.Context(), roomID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(messages)
}

func (h *MessageHandler) CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	roomID, err := h.chatService.InitiateChatRoom(c.Context(), req.CourseID, req.CustomerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"roomID": roomID,
	})
}
