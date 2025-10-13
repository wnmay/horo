package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/services/payment-service/internal/ports/outbound"
)

// Repository implements the PaymentRepository interface
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new PostgreSQL payment repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create saves a new payment to the database
func (r *Repository) Create(ctx context.Context, payment *domain.Payment) error {
	paymentModel := toPaymentModel(payment)
	result := r.db.WithContext(ctx).Create(paymentModel)
	return result.Error
}

// GetByID retrieves a payment by its ID
func (r *Repository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	var paymentModel Payment
	result := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).First(&paymentModel)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, result.Error
	}
	
	return toPaymentEntity(&paymentModel), nil
}

// GetByOrderID retrieves a payment by order ID
func (r *Repository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var paymentModel Payment
	result := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&paymentModel)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, result.Error
	}
	
	return toPaymentEntity(&paymentModel), nil
}

// Update saves changes to an existing payment
func (r *Repository) Update(ctx context.Context, payment *domain.Payment) error {
	paymentModel := toPaymentModel(payment)
	result := r.db.WithContext(ctx).Save(paymentModel)
	return result.Error
}

// Delete removes a payment from the database
func (r *Repository) Delete(ctx context.Context, paymentID string) error {
	result := r.db.WithContext(ctx).Delete(&Payment{}, "payment_id = ?", paymentID)
	return result.Error
}

// AutoMigrate runs database migrations for the Payment table
func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Payment{})
}

// Mapping functions between domain entity and database model
func toPaymentModel(payment *domain.Payment) *Payment {
	return &Payment{
		PaymentID:     payment.PaymentID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentStatus: PaymentStatus(payment.PaymentStatus),
		PaymentDate:   payment.PaymentDate,
	}
}

func toPaymentEntity(model *Payment) *domain.Payment {
	return &domain.Payment{
		PaymentID:     model.PaymentID,
		OrderID:       model.OrderID,
		Amount:        model.Amount,
		PaymentStatus: domain.PaymentStatus(model.PaymentStatus),
		PaymentDate:   model.PaymentDate,
	}
}