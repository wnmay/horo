package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	service "github.com/wnmay/horo/services/chat-service/internal/app"
	"github.com/wnmay/horo/services/chat-service/internal/domain"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type orderCompletedConsumer struct {
	chatService inbound_port.ChatService
	rmq         *message.RabbitMQ
}

func NewOrderCompletedConsumer(chatService inbound_port.ChatService, rmq *message.RabbitMQ) inbound_port.MessageConsumer {
	return &orderCompletedConsumer{
		chatService: chatService,
		rmq:         rmq,
	}
}

func (c *orderCompletedConsumer) StartListening() error {
	return c.rmq.ConsumeMessages(message.NotifyOrderCompleted, c.handleOrderCompletedCreated)
}

func (c *orderCompletedConsumer) handleOrderCompletedCreated(ctx context.Context, delivery amqp.Delivery) error {
	var amqpMessage contract.AmqpMessage
	var orderCompletedData message.OrderCompletedData

	// Parse the AMQP message
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Fatalf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	if err := json.Unmarshal(amqpMessage.Data, &orderCompletedData); err != nil {
		log.Printf("Failed to unmarshal message data: %v", err)
		return err
	}
	content := service.GenerateOrderCompletedMessage(orderCompletedData.OrderID, orderCompletedData.CourseID, orderCompletedData.OrderStatus, orderCompletedData.CourseName)

	messageID, err := c.chatService.SaveMessage(ctx, orderCompletedData.RoomID, "system", content)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return err
	}
	notificationData := message.ChatNotificationOutgoingData[message.OrderCompletedNotificationData]{
		MessageID: messageID,
		RoomID:    orderCompletedData.RoomID,
		SenderID:  "system",
		Type:      string(domain.MessageTypeNotification),
		CreatedAt: time.Now().Format(time.RFC3339),
		MessageDetail: &message.OrderCompletedNotificationData{
			OrderID:     orderCompletedData.OrderID,
			CourseID:    orderCompletedData.CourseID,
			OrderStatus: orderCompletedData.OrderStatus,
			CourseName:  orderCompletedData.CourseName,
		},
	}

	err = c.chatService.PublishOrderCompletedNotification(ctx, notificationData)
	if err != nil {
		log.Printf("Failed to publish orderCompleted created message: %v", err)
		return err
	}

	log.Printf("Published orderCompleted created message: %s", messageID)
	return nil
}
