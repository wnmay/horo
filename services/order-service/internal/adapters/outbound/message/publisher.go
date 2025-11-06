package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/wnmay/horo/services/order-service/internal/adapters/outbound/grpc"
	"github.com/wnmay/horo/services/order-service/internal/domain"
	"github.com/wnmay/horo/services/order-service/internal/ports/outbound"
	"github.com/wnmay/horo/shared/contract"
	"github.com/wnmay/horo/shared/message"
)

type Publisher struct {
	rabbit       *message.RabbitMQ
	courseClient *grpc.CourseClient
}

func NewPublisher(rabbit *message.RabbitMQ, courseClient *grpc.CourseClient) outbound.EventPublisher {
	return &Publisher{
		rabbit:       rabbit,
		courseClient: courseClient,
	}
}

func (p *Publisher) PublishOrderCreated(ctx context.Context, order *domain.Order) error {
	course, err := p.courseClient.GetCourseByID(ctx, order.CourseID)

	// Create order data for the event
	orderData := message.OrderData{
		OrderID:    order.OrderID.String(),
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		Amount:     course.Price,
	}

	// Marshal the order data
	data, err := json.Marshal(orderData)
	if err != nil {
		return fmt.Errorf("failed to marshal order data: %w", err)
	}

	// Create AMQP message using contract structure
	amqpMessage := contract.AmqpMessage{
		OwnerID: order.OrderID.String(),
		Data:    data,
	}

	// Publish the message with routing key using PublishMessage method
	if err := p.rabbit.PublishMessage(ctx, contract.OrderCreatedEvent, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish order created event: %w", err)
	}

	fmt.Printf("Published order created event for order: %s with price: %.2f\n", order.OrderID, course.Price)
	return nil
}

func (p *Publisher) PublishOrderCompleted(ctx context.Context, order *domain.Order) error {
	// Fetch course details from course service
	course, err := p.courseClient.GetCourseByID(ctx, order.CourseID)
	if err != nil {
		log.Printf("Warning: Failed to fetch course details for %s: %v. Using default course name.", order.CourseID, err)
		// Continue with default course name if fetch fails
	}

	orderCompletedData := message.OrderCompletedData{
		OrderID:     order.OrderID.String(),
		CourseID:    order.CourseID,
		CourseName:  course.Coursename,
		OrderStatus: string(order.Status),
		ProphetID:   course.ProphetId,
	}

	data, err := json.Marshal(orderCompletedData)
	if err != nil {
		return fmt.Errorf("failed to marshal order completion data: %w", err)
	}

	amqpMessage := contract.AmqpMessage{
		OwnerID: order.OrderID.String(),
		Data:    data,
	}

	if err := p.rabbit.PublishMessage(ctx, contract.OrderCompletedEvent, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish order completed event: %w", err)
	}

	fmt.Printf("Published order completed event for order: %s, course: %s, prophet: %s\n", order.OrderID, course.Coursename, course.ProphetId)
	return nil
}

func (p *Publisher) PublishOrderPaid(ctx context.Context, order *domain.Order) error {
	course, err := p.courseClient.GetCourseByID(ctx, order.CourseID)
	orderPaidData := message.OrderPaidData{
		OrderID:       order.OrderID.String(),
		PaymentID:     order.PaymentID.String(),
		RoomID:        order.RoomID.String(),
		CustomerID:    order.CustomerID.String(),
		CourseID:      order.CourseID.String(),
		OrderStatus:   string(order.Status),
		CourseName:    course.Coursename,
		Amount:       course.Price,
		PaymentStatus: "COMPLETED",
	}

	data, err := json.Marshal(orderPaidData)
	if err != nil {
		return fmt.Errorf("failed to marshal order paid data: %w", err)
	}

	amqpMessage := contract.AmqpMessage{
		OwnerID: order.OrderID.String(),
		Data:    data,
	}

	if err := p.rabbit.PublishMessage(ctx, contract.OrderPaidEvent, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish order paid event: %w", err)
	}

	fmt.Printf("Published order paid event for order: %s, payment: %s\n", order.OrderID, order.PaymentID)
	return nil
}

func (p *Publisher) PublishOrderPaymentBound(ctx context.Context, order *domain.Order) error {
	course, err := p.courseClient.GetCourseByID(ctx, order.CourseID)
	orderPaymentBoundData := message.OrderPaymentBoundData{
		OrderID:      order.OrderID.String(),
		PaymentID:    order.PaymentID.String(),
		RoomID:       order.RoomID.String(),
		CustomerID:   order.CustomerID.String(),
		OrderStatus:  string(order.Status),
		CourseID:     order.CourseID.String(),
		CourseName:   course.Coursename,
		Amount:      course.Price,
		PaymentStatus: "PENDING",
	}

	data, err := json.Marshal(orderPaymentBoundData)
	if err != nil {
		return fmt.Errorf("failed to marshal order payment bound data: %w", err)
	}

	amqpMessage := contract.AmqpMessage{
		OwnerID: order.OrderID.String(),
		Data:    data,
	}

	if err := p.rabbit.PublishMessage(ctx, contract.OrderPaymentBoundEvent, amqpMessage); err != nil {
		return fmt.Errorf("failed to publish order payment bound event: %w", err)
	}

	fmt.Printf("Published order payment bound event for order: %s, payment: %s\n", order.OrderID, order.PaymentID)
	return nil
}
