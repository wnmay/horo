package app

import (
	"context"
	"fmt"
	"log"
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

	if err := s.eventPublisher.PublishPaymentCreated(ctx, payment); err != nil {
		return nil, fmt.Errorf("payment created but failed to publish success event: %w", err)
	}

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

	payment.Status = status
	payment.UpdatedAt = time.Now()

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

func (s *Service) SettlePayment(ctx context.Context, orderID string) error {
	payment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	prev := payment.Status
	if err := payment.Settle(); err != nil {
		return fmt.Errorf("failed to settle payment: %w", err)
	}

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	if prev != domain.PaymentStatusSettled && payment.Status == domain.PaymentStatusSettled {
		if err := s.eventPublisher.PublishPaymentSettled(ctx, payment); err != nil {
			log.Printf("Payment settled but failed to publish event: %v", err)
		}
	}

	log.Printf("Payment %s settled successfully", payment.PaymentID)
	return nil
}
