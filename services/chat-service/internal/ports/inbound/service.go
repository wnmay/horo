package inbound_port

import (
	"context"
)

type ChatService interface {
	SaveMessage(ctx context.Context, roomID, senderID, content string) error
	InitiateChatRoom(ctx context.Context, courseID string, customerID string) (string, error)
}
