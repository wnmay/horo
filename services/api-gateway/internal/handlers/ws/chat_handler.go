package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/wnmay/horo/services/api-gateway/internal/messaging/publishers"
	gwWS "github.com/wnmay/horo/services/api-gateway/internal/websocket"
	chatpb "github.com/wnmay/horo/shared/proto/chat"
)

type ChatWSHandler struct {
	hub       *gwWS.Hub
	publisher *publishers.ChatMessagePublisher
	chatClient chatpb.ChatServiceClient 
}

func NewChatWSHandler(hub *gwWS.Hub, pub *publishers.ChatMessagePublisher,chatClient chatpb.ChatServiceClient) *ChatWSHandler {
	return &ChatWSHandler{hub: hub, publisher: pub,chatClient: chatClient}
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

	conn := gwWS.NewConnection(c)
	h.hub.AddUser(userID, conn)
	defer func() {
		h.hub.Remove(conn)
		conn.Close()
		log.Printf("[ws] disconnected user=%s", userID)
	}()

	initialRoomID := c.Query("roomId")
	if initialRoomID != "" {
		if h.validateAndJoinRoom(initialRoomID, userID, conn) {
			log.Printf("[ws] user=%s joined initial room=%s", userID, initialRoomID)
		}
	}

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
			if h.validateAndJoinRoom(msg.RoomID, userID, conn) {
				log.Printf("[ws] user=%s joined room=%s", userID, msg.RoomID)
				h.sendAck(conn, msg.RoomID)
			}

		case "message":
			if msg.RoomID == "" {
				log.Println("[ws] missing roomId in message")
				continue
			}
			if err := h.publisher.PublishMessageIncoming(
				context.Background(),
				msg.RoomID,
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

func (h *ChatWSHandler) validateAndJoinRoom(roomID, userID string, conn *gwWS.Connection) bool {
	resp, err := h.chatClient.ValidateRoomAccess(context.Background(), &chatpb.ValidateRoomRequest{
		UserId: userID,
		RoomId: roomID,
	})
	if err != nil {
		log.Printf("[ws] gRPC validation error for user=%s room=%s: %v", userID, roomID, err)
		_ = conn.Send([]byte(`{"error":"validation failed"}`))
		return false
	}
	if !resp.Allowed {
		log.Printf("[ws] user=%s denied joining room=%s: %s", userID, roomID, resp.Reason)
		_ = conn.Send([]byte(`{"error":"` + resp.Reason + `"}`))
		return false
	}

	h.hub.AddToRoom(roomID, conn)
	return true
}

func (h *ChatWSHandler) sendAck(conn *gwWS.Connection, roomID string) {
	ack := map[string]any{
		"type":   "joined",
		"roomId": roomID,
	}
	if ackBytes, err := json.Marshal(ack); err == nil {
		_ = conn.Send(ackBytes)
	}
}