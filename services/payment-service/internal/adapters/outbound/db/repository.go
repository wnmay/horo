package db

import (
	"context"
	"log"
	"time"
	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/services/payment-service/internal/ports/outbound"
	"gorm.io/gorm"
)

type paymentModel struct {
	PaymentID string    `gorm:"primaryKey;type:uuid;column:payment_id"`
	OrderID   string    `gorm:"not null;index;type:uuid"`
	Amount    float64   `gorm:"not null"`
	Status    string    `gorm:"not null;default:pending"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (paymentModel) TableName() string { return "payments" }

type GormPaymentRepository struct{ db *gorm.DB }

var _ outbound.PaymentRepository = (*GormPaymentRepository)(nil)

func NewGormPaymentRepository(db *gorm.DB) *GormPaymentRepository {
	// Auto-migrate payment table
	if err := db.AutoMigrate(&paymentModel{}); err != nil {
		log.Printf("Payment migration failed: %v", err)
	} else {
		log.Printf("Payments table migrated successfully")
	}
	
	return &GormPaymentRepository{db: db}
}

func (r *GormPaymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	model := paymentModel{
		PaymentID: p.PaymentID,
		OrderID:   p.OrderID,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *GormPaymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	var model paymentModel
	if err := r.db.WithContext(ctx).First(&model, "payment_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Payment{
		PaymentID: model.PaymentID,
		OrderID:   model.OrderID,
		Amount:    model.Amount,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *GormPaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var model paymentModel
	if err := r.db.WithContext(ctx).First(&model, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &domain.Payment{
		PaymentID: model.PaymentID,
		OrderID:   model.OrderID,
		Amount:    model.Amount,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *GormPaymentRepository) Update(ctx context.Context, p *domain.Payment) error {
	model := paymentModel{
		PaymentID: p.PaymentID,
		OrderID:   p.OrderID,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *GormPaymentRepository) Delete(ctx context.Context, paymentID string) error {
	return r.db.WithContext(ctx).Delete(&paymentModel{}, "payment_id = ?", paymentID).Error
}