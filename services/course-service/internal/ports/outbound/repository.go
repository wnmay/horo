package outbound

import "github.com/wnmay/horo/services/course-service/internal/domain"

type CourseRepository interface {
	Save(course *domain.Course) error
	FindByID(id string) (*domain.Course, error)
	FindAllByProphet(prophetID string) ([]*domain.Course, error)
}
