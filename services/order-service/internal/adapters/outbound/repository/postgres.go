package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/wnmay/horo/services/order-service/internal/domain/entity"
)

// Repository implements the OrderRepository interface
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new PostgreSQL order repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create saves a new order to the database
func (r *Repository) Create(ctx context.Context, order *entity.Order) error {
	orderModel := toOrderModel(order)
	result := r.db.WithContext(ctx).Create(orderModel)
	return result.Error
}

// GetByID retrieves an order by its ID
func (r *Repository) GetByID(ctx context.Context, orderID uuid.UUID) (*entity.Order, error) {
	var orderModel Order
	result := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&orderModel)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, result.Error
	}
	
	return toOrderEntity(&orderModel), nil
}

// GetByCustomerID retrieves all orders for a specific customer
func (r *Repository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entity.Order, error) {
	var orderModels []Order
	result := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&orderModels)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	orders := make([]*entity.Order, len(orderModels))
	for i, model := range orderModels {
		orders[i] = toOrderEntity(&model)
	}
	
	return orders, nil
}

// Update saves changes to an existing order
func (r *Repository) Update(ctx context.Context, order *entity.Order) error {
	orderModel := toOrderModel(order)
	result := r.db.WithContext(ctx).Save(orderModel)
	return result.Error
}

// Delete removes an order from the database
func (r *Repository) Delete(ctx context.Context, orderID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&Order{}, "order_id = ?", orderID)
	return result.Error
}

// AutoMigrate runs database migrations for the Order table
func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Order{})
}

// Mapping functions between domain entity and database model
func toOrderModel(order *entity.Order) *Order {
	model := &Order{
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		CourseID:   order.CourseID,
		Status:     OrderStatus(order.Status),
		Amount:     order.Amount,
		OrderDate:  order.OrderDate,
	}
	
	if order.PaymentID != nil {
		model.PaymentID = *order.PaymentID
	}
	
	return model
}

func toOrderEntity(model *Order) *entity.Order {
	order := &entity.Order{
		OrderID:    model.OrderID,
		CustomerID: model.CustomerID,
		CourseID:   model.CourseID,
		Status:     entity.OrderStatus(model.Status),
		Amount:     model.Amount,
		OrderDate:  model.OrderDate,
	}
	
	if model.PaymentID != (uuid.UUID{}) {
		order.PaymentID = &model.PaymentID
	}
	
	return order
}
