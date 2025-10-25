package app

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"	
	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/services/payment-service/internal/ports/inbound"
	"github.com/wnmay/horo/services/payment-service/internal/ports/outbound"
)

type Service struct {
	paymentRepo    outbound.PaymentRepository
	eventPublisher outbound.PaymentEventPublisher
}

func NewPaymentService(paymentRepo outbound.PaymentRepository, eventPublisher outbound.PaymentEventPublisher) *Service {
	return &Service{
		paymentRepo:    paymentRepo,
		eventPublisher: eventPublisher,
	}
}

// Payment Service Implementation
func (s *Service) CreatePaymentFromOrder(ctx context.Context, cmd inbound.CreatePaymentCommand) (*domain.Payment, error) {
	log.Printf("Creating payment for order: %s, amount: %.2f", cmd.OrderID, cmd.Amount)

	// Create new payment entity
	payment := domain.NewPayment(cmd.OrderID, cmd.Amount)

	// Save payment to repository
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// TODO: Here you could integrate with actual payment processor
	// For now, we'll simulate payment processing
	go s.simulatePaymentProcessing(ctx, payment)

	return payment, nil
}

func (s *Service) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return payment, nil
}

func (s *Service) GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment for order: %w", err)
	}
	return payment, nil
}

func (s *Service) UpdatePaymentStatus(ctx context.Context, paymentID string, status domain.PaymentStatus) error {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	payment.PaymentStatus = status
	payment.PaymentDate = time.Now()

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

func (s *Service) CompletePayment(ctx context.Context, paymentID string) error {
	// Get payment
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Complete the payment
	payment.Complete()

	// Update payment in repository
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Publish payment completed event
	if err := s.eventPublisher.PublishPaymentCompleted(ctx, payment); err != nil {
		log.Printf("Failed to publish payment completed event: %v", err)
		// Don't return error as payment was successfully completed
	}

	log.Printf("Payment %s completed successfully", payment.PaymentID)
	return nil
}

func (s *Service) simulatePaymentProcessing(ctx context.Context, payment *domain.Payment) {
	// Simulate payment processing delay
	time.Sleep(2 * time.Second)

	// Simulate 90% success rate
	if rand.Float32() < 0.9 {
		payment.Complete()
		
		// Update payment in repository
		if err := s.paymentRepo.Update(ctx, payment); err != nil {
			log.Printf("Failed to update payment status: %v", err)
			return
		}

		// Publish payment success event to notify order service
		if err := s.eventPublisher.PublishPaymentCompleted(ctx, payment); err != nil {
			log.Printf("Failed to publish payment success event: %v", err)
		}
		
		log.Printf("Payment %s completed successfully", payment.PaymentID)
	} else {
		payment.Fail()
		
		// Update payment in repository
		if err := s.paymentRepo.Update(ctx, payment); err != nil {
			log.Printf("Failed to update payment status: %v", err)
			return
		}

		// Publish payment failure event
		if err := s.eventPublisher.PublishPaymentFailed(ctx, payment); err != nil {
			log.Printf("Failed to publish payment failure event: %v", err)
		}
		
		log.Printf("Payment %s failed", payment.PaymentID)
	}
}
