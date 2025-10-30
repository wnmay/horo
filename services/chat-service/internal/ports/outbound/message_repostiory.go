package outbound_port

import (
	"context"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *domain.Message) error
	FindMessagesByRoomID(ctx context.Context, roomID string) ([]*domain.Message, error)
}
