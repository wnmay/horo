package service

import (
	"context"
	"encoding/json"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	outbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/outbound"
	"github.com/wnmay/horo/shared/contract"
)

type chatService struct {
	messageRepo      outbound_port.MessageRepository
	roomRepo         outbound_port.RoomRepositoryPort
	messagePublisher outbound_port.MessagePublisher
}

type paymentCreatedChatMessage struct {
	RoomID      string  `json:"roomId"`
	SenderID    string  `json:"senderId"`
	Content     string  `json:"content"`
	PaymentID   string  `json:"paymentId"`
	OrderID     string  `json:"orderId"`
	Amount      float64 `json:"amount"`
	MessageType string  `json:"messageType"`
}

func NewChatService(messageRepo outbound_port.MessageRepository, roomRepo outbound_port.RoomRepositoryPort, messagePublisher outbound_port.MessagePublisher) inbound_port.ChatService {
	return &chatService{
		messageRepo:      messageRepo,
		roomRepo:         roomRepo,
		messagePublisher: messagePublisher,
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
	mockProphetID := "prophet-1234" // TO DO: Use real prophet ID from course service
	room := domain.CreateRoom(mockProphetID, courseID, customerID)

	roomID, err := s.roomRepo.CreateRoom(context.Background(), room)
	if err != nil {
		return "", err
	}
	return roomID, nil
}

func (s *chatService) PublishPaymentCreatedMessage(ctx context.Context, paymentID string, orderID string, status string, amount float64) error {
	message := domain.CreateMessage(
		"mock-room-id", // TO DO: filter rooomId from payment details after we enrich the payment event data
		"system",
		GeneratePaymentCreatedMessage(paymentID, orderID, status, amount),
		domain.MessageTypeNotification,
		domain.MessageStatusSent,
	)
	if err := s.messageRepo.SaveMessage(ctx, message); err != nil {
		return err
	}

	data, err := json.Marshal(paymentCreatedChatMessage{
		RoomID:      message.RoomID,
		SenderID:    message.SenderID,
		Content:     message.Content,
		PaymentID:   paymentID,
		OrderID:     orderID,
		MessageType: string(message.Type),
		Amount:      amount,
	})
	if err != nil {
		return err
	}

	return s.messagePublisher.Publish(ctx, contract.AmqpMessage{
		OwnerID: orderID,
		Data:    data,
	})
}

func (s *chatService) GetMessagesByRoomID(ctx context.Context, roomID string) ([]*domain.Message, error) {
	return s.messageRepo.FindMessagesByRoomID(ctx, roomID)
}

func (s *chatService) GetChatRoomsByCustomerID(ctx context.Context, customerID string) ([]*domain.Room, error) {
	return s.roomRepo.GetChatRoomsByCustomerID(ctx, customerID)
}

func (s *chatService) GetChatRoomsByProphetID(ctx context.Context, prophetID string) ([]*domain.Room, error) {
	return s.roomRepo.GetChatRoomsByProphetID(ctx, prophetID)
}