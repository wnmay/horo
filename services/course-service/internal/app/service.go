package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/course-service/internal/domain"
	"github.com/wnmay/horo/services/course-service/internal/ports/outbound"
)

type CourseService interface {
	CreateCourse(input CreateCourseInput) (*domain.Course, error)
	GetCourseByID(id string) (*domain.Course, error)
	ListCoursesByProphet(prophetID string) ([]*domain.Course, error)
}

type courseService struct {
	repo outbound.CourseRepository
}

func (s courseService) GetCourseByID(id string) (*domain.Course, error) {
	return s.repo.FindByID(id)
}

func (s courseService) ListCoursesByProphet(prophetID string) ([]*domain.Course, error) {
	return s.repo.FindAllByProphet(prophetID)
}

func NewCourseService(r outbound.CourseRepository) CourseService {
	return &courseService{repo: r}
}

func (s *courseService) CreateCourse(input CreateCourseInput) (*domain.Course, error) {
	c := &domain.Course{
		ID:          generateID(),
		ProphetID:   input.ProphetID,
		CourseName:  input.CourseName,
		Description: input.Description,
		Price:       input.Price,
		Duration:    input.Duration,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Save(c); err != nil {
		return nil, err
	}
	return c, nil
}

func generateID() string {
	return "COURSE-" + uuid.New().String()
}

type CreateCourseInput struct {
	ID          string
	ProphetID   string
	CourseName  string
	Description string
	Price       float64
	Duration    domain.DurationEnum
	CreatedAt   time.Time
}
