package inbound

import "github.com/wnmay/horo/services/course-service/internal/domain"

type CourseServicePort interface {
	CreateCourse(course domain.Course) (*domain.Course, error)
	GetCourseByID(id string) (*domain.Course, error)
	ListCoursesByProphet(prophetID string) ([]*domain.Course, error)
}
