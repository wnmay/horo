package http_handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
)

type ChatHandler struct {
	chatService inbound_port.ChatService
}

type CreateRoomRequest struct {
	CourseID string `json:"courseID"`
}


type ValidateRoomAccessRequest struct {
	RoomID string `json:"roomID"`
}

func NewMessageHandler(chatService inbound_port.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) GetMessagesByRoomID(c *fiber.Ctx) error {
	roomID := c.Params("roomID")
	messages, err := h.chatService.GetMessagesByRoomID(c.Context(), roomID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(messages)
}

func (h *ChatHandler) CreateRoom(c *fiber.Ctx) error {
	// Check if handler is properly initialized
	if h == nil || h.chatService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Handler not properly initialized",
		})
	}

	var req CreateRoomRequest
	customerID := c.Get("X-User-Id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Validate required fields
	if req.CourseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "courseID is required",
		})
	}

	if customerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "X-User-Id header is required",
		})
	}

	roomID, err := h.chatService.InitiateChatRoom(c.Context(), req.CourseID, customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Println("Chat room created successfully with ID:", roomID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"roomID": roomID,
	})

}

func (h *ChatHandler) GetChatRoomsByCustomerID(c *fiber.Ctx) error {
	customerID := c.Get("X-User-Id")
	rooms, err := h.chatService.GetChatRoomsByCustomerID(c.Context(), customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(rooms)
}

func (h *ChatHandler) GetChatRoomsByProphetID(c *fiber.Ctx) error {
	prophetID := c.Get("X-User-Id")
	rooms, err := h.chatService.GetChatRoomsByProphetID(c.Context(), prophetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(rooms)
}

func (h *ChatHandler) ValidateRoomAccess(c *fiber.Ctx) error {
	userID := c.Get("X-User-Id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "X-User-Id header is required",
		})
	}

	var req ValidateRoomAccessRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}
	if req.RoomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "roomID is required",
		})
	}

	allowed, reason, err := h.chatService.ValidateRoomAccess(c.Context(), userID, req.RoomID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"allowed": allowed,
		"reason":  reason,
	})
}
