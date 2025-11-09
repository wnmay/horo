package inbound

import "github.com/wnmay/horo/services/course-service/internal/domain"

type CourseServicePort interface {
	CreateCourse(course domain.Course) (*domain.Course, error)
	GetCourseByID(id string) (*domain.Course, error)
	ListCoursesByProphet(prophetID string) ([]*domain.Course, error)
	UpdateCourse(id string, input *domain.UpdateCourseInput) (*domain.Course, error)
	DeleteCourse(id string) error
	FindCoursesByFilter(filter map[string]interface{}) ([]*domain.Course, error)
	CreateReview(review domain.Review) (*domain.Review, error)
	GetReviewByID(id string) (*domain.Review, error)
	ListReviewsByCourse(courseId string) ([]*domain.Review, error)
}
