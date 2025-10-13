package db

import (
	"context"
	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/services/payment-service/internal/ports/outbound"
	"gorm.io/gorm"
)

type paymentModel struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	OrderID  string `gorm:"not null;index;type:uuid"`
	UserID   string `gorm:"not null;type:uuid"`
	Amount   float64 `gorm:"not null"`
	Currency string `gorm:"not null;size:3"`
	Status   string `gorm:"not null;default:pending"`
}

func (paymentModel) TableName() string { return "payments" }

type GormPaymentRepository struct{ db *gorm.DB }

var _ outbound.PaymentRepository = (*GormPaymentRepository)(nil)

func NewGormPaymentRepository(db *gorm.DB) *GormPaymentRepository {
	// Auto-migrate payment table
	_ = db.AutoMigrate(&paymentModel{})
	return &GormPaymentRepository{db: db}
}

func (r *GormPaymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	model := paymentModel{
		ID:       p.ID,
		OrderID:  p.OrderID,
		UserID:   p.UserID,
		Amount:   p.Amount,
		Currency: p.Currency,
		Status:   p.Status,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *GormPaymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	var model paymentModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Payment{
		ID:       model.ID,
		OrderID:  model.OrderID,
		UserID:   model.UserID,
		Amount:   model.Amount,
		Currency: model.Currency,
		Status:   model.Status,
	}, nil
}

func (r *GormPaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var model paymentModel
	if err := r.db.WithContext(ctx).First(&model, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &domain.Payment{
		ID:       model.ID,
		OrderID:  model.OrderID,
		UserID:   model.UserID,
		Amount:   model.Amount,
		Currency: model.Currency,
		Status:   model.Status,
	}, nil
}

func (r *GormPaymentRepository) Update(ctx context.Context, p *domain.Payment) error {
	model := paymentModel{
		ID:       p.ID,
		OrderID:  p.OrderID,
		UserID:   p.UserID,
		Amount:   p.Amount,
		Currency: p.Currency,
		Status:   p.Status,
	}
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *GormPaymentRepository) Delete(ctx context.Context, paymentID string) error {
	return r.db.WithContext(ctx).Delete(&paymentModel{}, "id = ?", paymentID).Error
}