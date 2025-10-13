package inbound

import "github.com/wnmay/horo/services/order-service/internal/domain"

type OrderService interface {
	Create(name string) (domain.Person, error)
	GetAll() ([]domain.Person, error)
}
