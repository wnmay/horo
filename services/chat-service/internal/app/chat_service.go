package service

import (
	"context"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	outbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/outbound"
)

type chatService struct {
	messageRepo outbound_port.MessageRepository
	roomRepo    outbound_port.RoomRepositoryPort
}

func NewChatService(messageRepo outbound_port.MessageRepository, roomRepo outbound_port.RoomRepositoryPort) inbound_port.ChatService {
	return &chatService{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
	}
}

func (s *chatService) SaveMessage(ctx context.Context, roomID string, senderID string, content string) error {
	message := domain.CreateMessage(roomID, senderID, content, domain.MessageTypeText, domain.MessageStatusSent)
	if err := s.messageRepo.SaveMessage(context.Background(), message); err != nil {
		return err
	}
	return nil
}

func (s *chatService) InitiateChatRoom(ctx context.Context, courseID string, customerID string) (string, error) {
	mockProphetID := "prophet-1234"
	room := domain.CreateRoom(mockProphetID, courseID, customerID)
	
	if err := s.roomRepo.CreateRoom(context.Background(), room); err != nil {
		return "", err
	}
	return room.ID.Hex(), nil
}
