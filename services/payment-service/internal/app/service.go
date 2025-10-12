package app

import (
	"errors"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/services/payment-service/internal/ports/outbound"
)

type Service struct {
	repo outbound.PersonRepository
}

func NewService(r outbound.PersonRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(name string) (domain.Person, error) {
	if name == "" {
		return domain.Person{}, errors.New("name is required")
	}
	person := domain.Person{ID: uuid.NewString(), Name: name}
	if err := s.repo.Save(person); err != nil {
		return domain.Person{}, err
	}
	return person, nil
}

func (s *Service) GetAll() ([]domain.Person, error) {
	return s.repo.GetAll()
}
