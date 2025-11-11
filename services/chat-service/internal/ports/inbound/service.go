package inbound_port

import (
	"context"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
	"github.com/wnmay/horo/shared/message"
)

type ChatService interface {
	SaveMessage(ctx context.Context, roomID, senderID, content string, messageType domain.MessageType, status domain.MessageStatus, trigger string) (string, error)
	InitiateChatRoom(ctx context.Context, courseID string, customerID string) (string, error)
	PublishPaymentCreatedMessage(ctx context.Context, paymentID string, orderID string, status string, amount float64) error
	GetMessagesByRoomID(ctx context.Context, roomID string) ([]*domain.Message, error)
	GetChatRoomsByCustomerID(ctx context.Context, customerID string) ([]*domain.Room, error)
	GetChatRoomsByProphetID(ctx context.Context, prophetID string) ([]*domain.Room, error)
	PublishOutgoingMessage(ctx context.Context, message *domain.Message) error
	ValidateRoomAccess(ctx context.Context, userID, roomID string) (allowed bool, reason string, err error)
	GetChatRoomsByUserID(ctx context.Context, userID string) ([]*domain.RoomWithName, error)
	PublishOrderCompletedNotification(ctx context.Context, notificationData message.ChatNotificationOutgoingData[message.OrderCompletedNotificationData]) error
	PublishOrderPaymentBoundNotification(ctx context.Context, notificationData message.ChatNotificationOutgoingData[message.OrderPaymentBoundNotificationData]) error
	PublishOrderPaidNotification(ctx context.Context, notificationData message.ChatNotificationOutgoingData[message.OrderPaidNotificationData]) error
	UpdateRoomIsDone(ctx context.Context, roomID string, isDone bool) error
}
