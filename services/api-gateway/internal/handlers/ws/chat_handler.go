package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/wnmay/horo/services/api-gateway/internal/messaging/publishers"
	gwWS "github.com/wnmay/horo/services/api-gateway/internal/websocket"
)

type ChatWSHandler struct {
	hub       *gwWS.Hub
	publisher *publishers.ChatMessagePublisher
}

func NewChatWSHandler(hub *gwWS.Hub, pub *publishers.ChatMessagePublisher) *ChatWSHandler {
	return &ChatWSHandler{hub: hub, publisher: pub}
}

func (h *ChatWSHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/ws/chat", websocket.New(h.handle))
}

type wsMessage struct {
	Action  string `json:"action"`
	RoomID  string `json:"roomId"`
	Type    string `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
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

	initialRoomID := c.Query("roomId")
	conn := gwWS.NewConnection(c)
	h.hub.AddUser(userID, conn)

	if initialRoomID != "" {
		h.hub.AddToRoom(initialRoomID, conn)
		log.Printf("[ws] user=%s joined initial room=%s", userID, initialRoomID)
	}

	defer func() {
		h.hub.Remove(conn)
		conn.Close()
		log.Printf("[ws] disconnected user=%s", userID)
	}()

	for {
		raw, err := conn.Read()
		if err != nil {
			break
		}

		var msg wsMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("[ws] invalid json:", err)
			continue
		}

		switch msg.Action {
		case "join_room":
			if msg.RoomID == "" {
				log.Println("[ws] missing roomId in join_room")
				continue
			}
			h.hub.AddToRoom(msg.RoomID, conn)
			log.Printf("[ws] user=%s joined room=%s", userID, msg.RoomID)

			ack := map[string]any{
				"type":   "joined",
				"roomId": msg.RoomID,
			}
			if ackBytes, err := json.Marshal(ack); err == nil {
				_ = conn.Send(ackBytes)
			}

		case "message":
			roomID := msg.RoomID
			if roomID == "" {
				log.Println("[ws] missing roomId in message")
				continue
			}
			if err := h.publisher.PublishMessageIncoming(
				context.Background(),
				roomID,
				userID,
				msg.Content,
				msg.Type,
			); err != nil {
				log.Printf("[ws] publish err: %v", err)
			}

		default:
			log.Printf("[ws] unknown action: %s", msg.Action)
		}
	}
}