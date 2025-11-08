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
	RoomExists(ctx context.Context, roomID string) (bool, error)
	IsUserInRoom(ctx context.Context, roomID string, userID string) (bool, error)
	GetChatRoomsByUserID(ctx context.Context, userID string) ([]*domain.Room, error)
}