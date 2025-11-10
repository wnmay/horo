package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	outbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/outbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
	shared_message "github.com/wnmay/horo/shared/message"
)

type chatService struct {
	messageRepo      outbound_port.MessageRepository
	roomRepo         outbound_port.RoomRepositoryPort
	messagePublisher outbound_port.MessagePublisher
	userProvider     outbound_port.UserProvider
}

type paymentCreatedChatMessage struct {
	MessageID   string  `json:"messageId"`
	RoomID      string  `json:"roomId"`
	SenderID    string  `json:"senderId"`
	Content     string  `json:"content"`
	PaymentID   string  `json:"paymentId"`
	OrderID     string  `json:"orderId"`
	Amount      float64 `json:"amount"`
	MessageType string  `json:"messageType"`
}

func NewChatService(messageRepo outbound_port.MessageRepository, roomRepo outbound_port.RoomRepositoryPort, messagePublisher outbound_port.MessagePublisher, userProvider outbound_port.UserProvider) inbound_port.ChatService {
	return &chatService{
		messageRepo:      messageRepo,
		roomRepo:         roomRepo,
		messagePublisher: messagePublisher,
		userProvider:     userProvider,
	}
}

func (s *chatService) SaveMessage(ctx context.Context, roomID string, senderID string, content string) (string, error) {
	message := domain.CreateMessage("", roomID, senderID, content, domain.MessageTypeText, domain.MessageStatusSent)
	messageID, err := s.messageRepo.SaveMessage(context.Background(), message)
	if err != nil {
		return "", err
	}
	return messageID, nil
}

func (s *chatService) InitiateChatRoom(ctx context.Context, courseID string, customerID string) (string, error) {
	mockProphetID := "prophet-1234" // TO DO: Use real prophet ID from course service
	room := domain.CreateRoom(mockProphetID, customerID, courseID, false)

	roomID, err := s.roomRepo.CreateRoom(context.Background(), room)
	if err != nil {
		log.Println("Error creating chat room:", err)
		return "", err
	}

	return roomID, nil
}

func (s *chatService) PublishPaymentCreatedMessage(ctx context.Context, paymentID string, orderID string, status string, amount float64) error {
	message := domain.CreateMessage(
		"",
		"mock-room-id", // TO DO: filter rooomId from payment details after we enrich the payment event data
		"system",
		GeneratePaymentCreatedMessage(paymentID, orderID, status, amount),
		domain.MessageTypeNotification,
		domain.MessageStatusSent,
	)
	messageID, err := s.messageRepo.SaveMessage(ctx, message)
	if err != nil {
		return err
	}

	data, err := json.Marshal(paymentCreatedChatMessage{
		MessageID:   messageID,
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

// Publish message from another user to the chat room
func (s *chatService) PublishOutgoingMessage(ctx context.Context, message *domain.Message) error {
	messageData := shared_message.ChatMessageOutgoingData{
		MessageID: message.ID,
		RoomID:    message.RoomID,
		SenderID:  message.SenderID,
		Content:   message.Content,
		Type:      string(message.Type),
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
	}
	data, err := json.Marshal(messageData)
	if err != nil {
		return err
	}
	return s.messagePublisher.Publish(ctx, contract.AmqpMessage{
		OwnerID: message.SenderID,
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

func (s *chatService) ValidateRoomAccess(ctx context.Context, userID string, roomID string) (bool, string, error) {
	exists, err := s.roomRepo.RoomExists(ctx, roomID)
	if err != nil {
		return false, "internal error", err
	}
	if !exists {
		return false, "room not found", nil
	}
	joinable, err := s.roomRepo.IsUserInRoom(ctx, userID, roomID)
	if err != nil {
		return false, "internal error", err
	}
	if !joinable {
		return false, "user cannot chat in this room", nil
	}

	return true, "", nil
}

func (s *chatService) GetChatRoomsByUserID(ctx context.Context, userID string) ([]*domain.RoomWithName, error) {
	rooms, err := s.roomRepo.GetChatRoomsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var userIDs []string
	for _, r := range rooms {
		userIDs = append(userIDs, r.CustomerID)
		userIDs = append(userIDs, r.ProphetID)
	}
	users, err := s.userProvider.MapUserNamesByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	var roomWithNames []*domain.RoomWithName

	for _, r := range rooms {
		roomWithNames = append(roomWithNames, &domain.RoomWithName{
			ID:           r.ID,
			ProphetID:    r.ProphetID,
			CustomerID:   r.CustomerID,
			CourseID:     r.CourseID,
			CreatedAt:    r.CreatedAt,
			LastMessage:  r.LastMessage,
			IsDone:       r.IsDone,
			ProphetName:  users[r.ProphetID].Name,
			CustomerName: users[r.CustomerID].Name,
		})
	}
	return roomWithNames, nil

}

func (s *chatService) PublishOrderCompletedNotification(ctx context.Context, notificationData message.ChatNotificationOutgoingData[message.OrderCompletedNotificationData]) error {
	data, err := json.Marshal(notificationData)
	if err != nil {
		return err
	}

	return s.messagePublisher.Publish(ctx, contract.AmqpMessage{
		OwnerID: notificationData.SenderID,
		Data:    data,
	})
}

func (s *chatService) PublishOrderPaymentBoundNotification(ctx context.Context, notificationData message.ChatNotificationOutgoingData[message.OrderPaymentBoundNotificationData]) error {
	data, err := json.Marshal(notificationData)
	if err != nil {
		return err
	}
	return s.messagePublisher.Publish(ctx, contract.AmqpMessage{
		OwnerID: notificationData.SenderID,
		Data:    data,
	})
}

func (s *chatService) PublishOrderPaidNotification(ctx context.Context, notificationData message.ChatNotificationOutgoingData[message.OrderPaidNotificationData]) error {
	data, err := json.Marshal(notificationData)
	if err != nil {
		return err
	}
	return s.messagePublisher.Publish(ctx, contract.AmqpMessage{
		OwnerID: notificationData.SenderID,
		Data:    data,
	})
}

func (s *chatService) UpdateRoomIsDone(ctx context.Context, roomID string, isDone bool) error {
	return s.roomRepo.UpdateRoomIsDoneByRoomID(ctx, roomID, isDone)
}
