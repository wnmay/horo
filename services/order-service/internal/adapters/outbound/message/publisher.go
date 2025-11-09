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
	
	// Default values if course fetch fails
	coursePrice := 200.0
	if err != nil {
		log.Printf("Warning: Failed to fetch course details for %s: %v. Using default price.", order.CourseID, err)
	} else if course != nil {
		coursePrice = course.Price
	}

	// Create order data for the event
	orderData := message.OrderData{
		OrderID:    order.OrderID.String(),
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		Amount:     coursePrice,
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

	fmt.Printf("Published order created event for order: %s with price: %.2f\n", order.OrderID, coursePrice)
	return nil
}

func (p *Publisher) PublishOrderCompleted(ctx context.Context, order *domain.Order) error {
	// Fetch course details from course service
	course, err := p.courseClient.GetCourseByID(ctx, order.CourseID)
	
	// Default values if course fetch fails
	courseName := "Unknown Course"
	prophetID := ""
	
	if err != nil {
		log.Printf("Warning: Failed to fetch course details for %s: %v. Using default course name.", order.CourseID, err)
	} else if course != nil {
		courseName = course.Coursename
		prophetID = course.ProphetId
	}

	orderCompletedData := message.OrderCompletedData{
		OrderID:     order.OrderID.String(),
		CourseID:    order.CourseID,
		CourseName:  courseName,
		OrderStatus: string(order.Status),
		ProphetID:   prophetID,
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

	fmt.Printf("Published order completed event for order: %s, course: %s, prophet: %s\n", order.OrderID, courseName, prophetID)
	return nil
}

func (p *Publisher) PublishOrderPaid(ctx context.Context, order *domain.Order) error {
	course, err := p.courseClient.GetCourseByID(ctx, order.CourseID)
	
	// Default values if course fetch fails
	courseName := "Unknown Course"
	coursePrice := 200.0
	
	if err != nil {
		log.Printf("Warning: Failed to fetch course details for %s: %v. Using defaults.", order.CourseID, err)
	} else if course != nil {
		courseName = course.Coursename
		coursePrice = course.Price
	}
	
	orderPaidData := message.OrderPaidData{
		OrderID:       order.OrderID.String(),
		PaymentID:     order.PaymentID.String(),
		RoomID:        order.RoomID,
		CustomerID:    order.CustomerID,
		CourseID:      order.CourseID,
		OrderStatus:   string(order.Status),
		CourseName:    courseName,
		Amount:        coursePrice,
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
	
	// Default values if course fetch fails
	courseName := "Unknown Course"
	coursePrice := 200.0
	
	if err != nil {
		log.Printf("Warning: Failed to fetch course details for %s: %v. Using defaults.", order.CourseID, err)
	} else if course != nil {
		courseName = course.Coursename
		coursePrice = course.Price
	}
	
	orderPaymentBoundData := message.OrderPaymentBoundData{
		OrderID:       order.OrderID.String(),
		PaymentID:     order.PaymentID.String(),
		RoomID:        order.RoomID,
		CustomerID:    order.CustomerID,
		OrderStatus:   string(order.Status),
		CourseID:      order.CourseID,
		CourseName:    courseName,
		Amount:        coursePrice,
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
