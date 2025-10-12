package outbound

import "github.com/wnmay/horo/services/payment-service/internal/domain"

type PersonRepository interface {
	Save(p domain.Person) error
	GetAll() ([]domain.Person, error)
}