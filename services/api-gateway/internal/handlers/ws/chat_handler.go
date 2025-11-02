package ws

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	gwWS "github.com/wnmay/horo/services/api-gateway/internal/websocket"
)

type ChatWSHandler struct {
	hub *gwWS.Hub
}

func NewChatWSHandler(hub *gwWS.Hub) *ChatWSHandler {
	return &ChatWSHandler{hub: hub}
}

func (h *ChatWSHandler) RegisterRoutes(app *fiber.App) {
	app.Use("/ws/chat", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/chat", websocket.New(h.handle))
}

func (h *ChatWSHandler) handle(c *websocket.Conn) {
	userID, _ := c.Locals("userId").(string)
	if userID == "" {
		log.Println("unauthorized websocket attempt, closing")
		_ = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"),
		)
		_ = c.Close()
		return
	}

	conn := gwWS.NewConnection(c)
	h.hub.Add(userID, conn)
	log.Printf("[ws] connected")

	defer func() {
		h.hub.Remove(userID, conn)
		conn.Close()
		log.Printf("[ws] disconnected")
	}()

	for {
		msg, err := conn.Read()
		if err != nil {
			break
		}

		h.hub.SendTo(userID, msg)
	}
}
