package outbound_port

import (
	"context"
	"github.com/wnmay/horo/services/chat-service/internal/domain"
)
type RoomRepositoryPort interface {
	FindRoomByID(ctx context.Context, roomID string) (*domain.Room, error)
	CreateRoom(ctx context.Context, room *domain.Room) (string, error)
	GetChatRoomsByCustomerID(ctx context.Context, customerID string) ([]*domain.Room, error)
	GetChatRoomsByProphetID(ctx context.Context, prophetID string) ([]*domain.Room, error)
}