package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	service "github.com/wnmay/horo/services/chat-service/internal/app"
	"github.com/wnmay/horo/services/chat-service/internal/domain"
	inbound_port "github.com/wnmay/horo/services/chat-service/internal/ports/inbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type notificationConsumer struct {
	chatService inbound_port.ChatService
	rmq         *message.RabbitMQ
}

func NewNotificationConsumer(chatService inbound_port.ChatService, rmq *message.RabbitMQ) inbound_port.MessageConsumer {
	return &notificationConsumer{
		chatService: chatService,
		rmq:         rmq,
	}
}

func (c *notificationConsumer) StartListening() error {
	return c.rmq.ConsumeMessages(message.NotifyOrderCompleted, c.handleNotification)
}

// handleNotification routes messages to the appropriate handler based on routing key
func (c *notificationConsumer) handleNotification(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("Received notification event with routing key '%s': %s", delivery.RoutingKey, delivery.Body)

	// Route to appropriate handler based on routing key
	switch delivery.RoutingKey {
	case contract.OrderCompletedEvent:
		return c.handleOrderCompleted(ctx, delivery)
	case contract.OrderPaymentBoundEvent:
		return c.handleOrderPaymentBound(ctx, delivery)
	case contract.OrderPaidEvent:
		return c.handleOrderPaid(ctx, delivery)
	default:
		log.Printf("Unknown routing key: %s, skipping message", delivery.RoutingKey)
		return fmt.Errorf("unknown routing key: %s", delivery.RoutingKey)
	}
}

func (c *notificationConsumer) handleOrderCompleted(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("Handling order completed event")
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

	// Handle case where roomID might be empty - try to create room if needed
	roomID := orderCompletedData.RoomID
	if roomID == "" {
		log.Printf("RoomID is empty for order %s, attempting to create chat room", orderCompletedData.OrderID)
		// Create a new chat room
		newRoomID, err := c.chatService.InitiateChatRoom(ctx, orderCompletedData.CourseID, orderCompletedData.CustomerID)
		if err != nil {
			log.Printf("Failed to create chat room: %v", err)
			return err
		}
		roomID = newRoomID
		log.Printf("Created new chat room with ID: %s", roomID)
	}

	if err := c.chatService.UpdateRoomIsDone(ctx, roomID, true); err != nil {
		log.Printf("Failed to update room is done: %v", err)
		return err
	}

	content := service.GenerateOrderCompletedMessage(orderCompletedData.OrderID, orderCompletedData.CourseID, orderCompletedData.OrderStatus, orderCompletedData.CourseName)

	messageID, err := c.chatService.SaveMessage(ctx, roomID, "system", content, domain.MessageTypeNotification, domain.MessageStatusSent, string(contract.OrderCompletedEvent))
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return err
	}
	notificationData := message.ChatNotificationOutgoingData[message.OrderCompletedNotificationData]{
		MessageID: messageID,
		RoomID:    roomID,
		SenderID:  "system",
		Type:      string(domain.MessageTypeNotification),
		CreatedAt: time.Now().Format(time.RFC3339),
		MessageDetail: &message.OrderCompletedNotificationData{
			OrderID:     orderCompletedData.OrderID,
			CourseID:    orderCompletedData.CourseID,
			OrderStatus: orderCompletedData.OrderStatus,
			CourseName:  orderCompletedData.CourseName,
		},
		Trigger: contract.OrderCompletedEvent,
	}

	err = c.chatService.PublishOrderCompletedNotification(ctx, notificationData)
	if err != nil {
		log.Printf("Failed to publish orderCompleted created message: %v", err)
		return err
	}
	log.Printf("Published orderCompleted created message: %s", messageID)
	return nil
}

func (c *notificationConsumer) handleOrderPaymentBound(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("Handling order payment bound event")
	var amqpMessage contract.AmqpMessage
	var orderPaymentBoundData message.OrderPaymentBoundData

	// Parse the AMQP message
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Fatalf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	if err := json.Unmarshal(amqpMessage.Data, &orderPaymentBoundData); err != nil {
		log.Printf("Failed to unmarshal message data: %v", err)
		return err
	}

	// Generate notification content for payment bound
	content := service.GenerateOrderPaymentBoundMessage(orderPaymentBoundData.OrderID, orderPaymentBoundData.CourseID, orderPaymentBoundData.OrderStatus, orderPaymentBoundData.CourseName)

	messageID, err := c.chatService.SaveMessage(ctx, orderPaymentBoundData.RoomID, "system", content, domain.MessageTypeNotification, domain.MessageStatusSent, string(contract.OrderPaymentBoundEvent))
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return err
	}

	notificationData := message.ChatNotificationOutgoingData[message.OrderPaymentBoundNotificationData]{
		MessageID: messageID,
		RoomID:    orderPaymentBoundData.RoomID,
		SenderID:  "system",
		Type:      string(domain.MessageTypeNotification),
		CreatedAt: time.Now().Format(time.RFC3339),
		Trigger:   contract.OrderPaymentBoundEvent,
		MessageDetail: &message.OrderPaymentBoundNotificationData{
			OrderID:       orderPaymentBoundData.OrderID,
			PaymentID:     orderPaymentBoundData.PaymentID,
			RoomID:        orderPaymentBoundData.RoomID,
			CustomerID:    orderPaymentBoundData.CustomerID,
			CourseID:      orderPaymentBoundData.CourseID,
			OrderStatus:   orderPaymentBoundData.OrderStatus,
			CourseName:    orderPaymentBoundData.CourseName,
			Amount:        orderPaymentBoundData.Amount,
			PaymentStatus: orderPaymentBoundData.PaymentStatus,
		},
	}

	err = c.chatService.PublishOrderPaymentBoundNotification(ctx, notificationData)
	if err != nil {
		log.Printf("Failed to publish order payment bound notification: %v", err)
		return err
	}
	log.Printf("Published order payment bound notification: %s", messageID)
	return nil
}

func (c *notificationConsumer) handleOrderPaid(ctx context.Context, delivery amqp.Delivery) error {
	log.Printf("Handling order paid event")
	var amqpMessage contract.AmqpMessage
	var orderPaidData message.OrderPaidData

	// Parse the AMQP message
	if err := json.Unmarshal(delivery.Body, &amqpMessage); err != nil {
		log.Fatalf("Failed to unmarshal AMQP message: %v", err)
		return err
	}

	if err := json.Unmarshal(amqpMessage.Data, &orderPaidData); err != nil {
		log.Printf("Failed to unmarshal message data: %v", err)
		return err
	}

	// Generate notification content for payment success
	content := service.GenerateOrderPaidMessage(orderPaidData.OrderID, orderPaidData.CourseID, orderPaidData.OrderStatus, orderPaidData.CourseName)

	messageID, err := c.chatService.SaveMessage(ctx, orderPaidData.RoomID, "system", content, domain.MessageTypeNotification, domain.MessageStatusSent, string(contract.OrderPaidEvent))
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return err
	}

	notificationData := message.ChatNotificationOutgoingData[message.OrderPaidNotificationData]{
		MessageID: messageID,
		RoomID:    orderPaidData.RoomID,
		SenderID:  "system",
		Type:      string(domain.MessageTypeNotification),
		CreatedAt: time.Now().Format(time.RFC3339),
		Trigger:   contract.OrderPaidEvent,
		MessageDetail: &message.OrderPaidNotificationData{
			OrderID:       orderPaidData.OrderID,
			PaymentID:     orderPaidData.PaymentID,
			RoomID:        orderPaidData.RoomID,
			CustomerID:    orderPaidData.CustomerID,
			CourseID:      orderPaidData.CourseID,
			OrderStatus:   orderPaidData.OrderStatus,
			CourseName:    orderPaidData.CourseName,
			Amount:        orderPaidData.Amount,
			PaymentStatus: orderPaidData.PaymentStatus,
		},
	}

	err = c.chatService.PublishOrderPaidNotification(ctx, notificationData)
	if err != nil {
		log.Printf("Failed to publish order paid notification: %v", err)
		return err
	}
	log.Printf("Published order paid notification: %s", messageID)
	return nil
}
