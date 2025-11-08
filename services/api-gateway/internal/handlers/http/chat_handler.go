package http_handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/shared/env"
)

type ChatHandler struct {
	chatServiceURL string
	client         *http.Client
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		chatServiceURL: env.GetString("CHAT_SERVICE_URL", "http://localhost:3004"),
		client:         &http.Client{},
	}
}

func (h *ChatHandler) GetMessagesByRoomID(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.chatServiceURL, fmt.Sprintf("/api/chat/%s/messages", c.Params("roomID")))
}

func (h *ChatHandler) CreateRoom(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "POST", h.chatServiceURL, "/api/chat/rooms")
}

func (h *ChatHandler) GetChatRoomsByCustomerID(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.chatServiceURL, "/api/chat/customer/rooms")
}

func (h *ChatHandler) GetChatRoomsByProphetID(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.chatServiceURL, "/api/chat/prophet/rooms")
}

func (h *ChatHandler) GetChatRoomsByUserID(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.chatServiceURL, "/api/chat/user/rooms")
}