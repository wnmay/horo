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
	content := fmt.Sprintf("Payment created for order %s. Amount: %.2f, Status: %s",
		orderPaymentBoundData.OrderID,
		orderPaymentBoundData.Amount,
		orderPaymentBoundData.PaymentStatus)

	messageID, err := c.chatService.SaveMessage(ctx, orderPaymentBoundData.RoomID, "system", content)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return err
	}

	// You can add a publish notification here if needed
	log.Printf("Saved payment bound notification message: %s", messageID)
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
	content := fmt.Sprintf("Payment successful! Order %s for %s has been paid. Amount: %.2f",
		orderPaidData.OrderID,
		orderPaidData.CourseName,
		orderPaidData.Amount)

	messageID, err := c.chatService.SaveMessage(ctx, orderPaidData.RoomID, "system", content)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return err
	}

	// You can add a publish notification here if needed
	log.Printf("Saved order paid notification message: %s", messageID)
	return nil
}
