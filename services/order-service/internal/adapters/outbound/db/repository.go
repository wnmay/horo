package db

import (
	"time"

	"github.com/google/uuid"
	"context"
	"errors"
	"log"

	"github.com/wnmay/horo/services/order-service/internal/domain"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusCancelled  OrderStatus = "CANCELLED"
	StatusFailed     OrderStatus = "CONFIRMED"
)

type Order struct {
	OrderID     uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID  string      `gorm:"type:varchar(255);not null"`
	CourseID    uuid.UUID   `gorm:"type:uuid;not null"`
	PaymentID   uuid.UUID   `gorm:"type:uuid"` 
	Status      OrderStatus `gorm:"type:varchar(20);not null"`
	OrderDate   time.Time   `gorm:"not null"`
}

func (o *Order) TableName() string {
	return "orders"
}

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
func (r *Repository) Create(ctx context.Context, order *domain.Order) error {
	orderModel := toOrderModel(order)
	result := r.db.WithContext(ctx).Create(orderModel)
	return result.Error
}
// Get all orders
func (r *Repository) GetAll(ctx context.Context) ([]*domain.Order, error) {
	var orderModels []Order
	result := r.db.WithContext(ctx).Find(&orderModels)
	if result.Error != nil {
		return nil, result.Error
	}
	orders := make([]*domain.Order, len(orderModels))
	for i, model := range orderModels {
		orders[i] = toOrderEntity(&model)
	
	}
	return orders, nil
}
// GetByID retrieves an order by its ID
func (r *Repository) GetByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
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
func (r *Repository) GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Order, error) {
	var orderModels []Order
	result := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&orderModels)

	if result.Error != nil {
		return nil, result.Error
	}

	orders := make([]*domain.Order, len(orderModels))
	for i, model := range orderModels {
		orders[i] = toOrderEntity(&model)
	}

	return orders, nil
}

// Update saves changes to an existing order
func (r *Repository) Update(ctx context.Context, order *domain.Order) error {
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
	if err := r.db.AutoMigrate(&Order{}); err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Printf("Orders table migrated successfully")
	return nil
}

// Mapping functions between domain entity and database model
func toOrderModel(order *domain.Order) *Order {
	model := &Order{
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		CourseID:   order.CourseID,
		Status:     OrderStatus(order.Status),
		OrderDate:  order.OrderDate,
	}

	if order.PaymentID != nil {
		model.PaymentID = *order.PaymentID
	}

	return model
}

func toOrderEntity(model *Order) *domain.Order {
	order := &domain.Order{
		OrderID:    model.OrderID,
		CustomerID: model.CustomerID,
		CourseID:   model.CourseID,
		Status:     domain.OrderStatus(model.Status),
		OrderDate:  model.OrderDate,
	}

	if model.PaymentID != (uuid.UUID{}) {
		order.PaymentID = &model.PaymentID
	}

	return order
}
