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
	customerID := c.Get("X-User-Uid")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	roomID, err := h.chatService.InitiateChatRoom(c.Context(), req.CourseID, customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"roomID": roomID,
	})
}

func (h *MessageHandler) GetChatRoomsByCustomerID(c *fiber.Ctx) error {
	customerID := c.Get("X-User-Uid")
	rooms, err := h.chatService.GetChatRoomsByCustomerID(c.Context(), customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(rooms)
}

func (h *MessageHandler) GetChatRoomsByProphetID(c *fiber.Ctx) error {
	prophetID := c.Get("X-User-Uid")
	rooms, err := h.chatService.GetChatRoomsByProphetID(c.Context(), prophetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(rooms)
}