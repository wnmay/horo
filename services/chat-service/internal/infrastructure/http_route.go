package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	http_handler "github.com/wnmay/horo/services/chat-service/internal/adapters/inbound/http"
)

func SetupHTTPRoutes(app *fiber.App, ChatHandler *http_handler.ChatHandler) {
	api := app.Group("/api/chat")

	api.Get("/:roomID/messages", ChatHandler.GetMessagesByRoomID)
	api.Post("/rooms", ChatHandler.CreateRoom)
	api.Get("/customer/rooms", ChatHandler.GetChatRoomsByCustomerID)
	api.Get("/prophet/rooms", ChatHandler.GetChatRoomsByProphetID)
	api.Post("/room/validate", ChatHandler.ValidateRoomAccess)
}
