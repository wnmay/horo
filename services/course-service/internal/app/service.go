package app

import (
	"time"

	"github.com/wnmay/horo/services/course-service/internal/domain"
	"github.com/wnmay/horo/services/course-service/internal/ports/outbound"
)

type CourseService struct {
	repo outbound.CourseRepository
}

func NewCourseService(r outbound.CourseRepository) *CourseService {
	return &CourseService{repo: r}
}

func (s *CourseService) CreateCourse(input domain.Course) (*domain.Course, error) {
	input.ID = generateID()
	input.CreatedAt = time.Now()
	if err := s.repo.Save(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func generateID() any {
	panic("unimplemented")
}
