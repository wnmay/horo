package inbound_port

import (
	"context"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
)

type ChatService interface {
	SaveMessage(ctx context.Context, roomID, senderID, content string) error
	InitiateChatRoom(ctx context.Context, courseID string, customerID string) (string, error)
	PublishPaymentCreatedMessage(ctx context.Context, paymentID string, orderID string, status string, amount float64) error
	GetMessagesByRoomID(ctx context.Context, roomID string) ([]*domain.Message, error)
	GetChatRoomsByCustomerID(ctx context.Context, customerID string) ([]*domain.Room, error)
	GetChatRoomsByProphetID(ctx context.Context, prophetID string) ([]*domain.Room, error)
}
