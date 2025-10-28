package inbound_port

import (
	"context"
)

type ChatService interface {
	SaveMessage(ctx context.Context, roomID, senderID, content string) error
	InitiateChatRoom(ctx context.Context, courseID string, customerID string) (string, error)
	PublishPaymentCreatedMessage(ctx context.Context, paymentID string, orderID string, status string, amount float64) error
}
