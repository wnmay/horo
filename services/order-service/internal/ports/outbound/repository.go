package outbound

import "github.com/wnmay/horo/services/order-service/internal/domain"

type PersonRepository interface {
	Save(p domain.Person) error
	GetAll() ([]domain.Person, error)
}
