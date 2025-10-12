package inbound

import "github.com/wnmay/horo/services/payment-service/internal/domain"

type PersonService interface {
	Create(name string) (domain.Person, error)
	GetAll() ([]domain.Person, error)
}
