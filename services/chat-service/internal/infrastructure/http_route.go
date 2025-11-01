package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	http_handler "github.com/wnmay/horo/services/chat-service/internal/adapters/inbound/http"
)

func SetupHTTPRoutes(app *fiber.App, messageHandler *http_handler.MessageHandler) {
	api := app.Group("/api/chat")

	api.Get("/:roomID/messages", messageHandler.GetMessagesByRoomID)
	api.Post("/rooms", messageHandler.CreateRoom)
	api.Get("/customer/rooms", messageHandler.GetChatRoomsByCustomerID)
	api.Get("/prophet/rooms", messageHandler.GetChatRoomsByProphetID)
}
