package outbound

import "github.com/wnmay/horo/services/course-service/internal/domain"

type CourseRepository interface {
	SaveCourse(course *domain.Course) error
	FindCourseByID(id string) (*domain.Course, error)
	FindCoursesByProphet(prophetID string) ([]*domain.Course, error)
	Update(id string, updates map[string]interface{}) (*domain.Course, error)
	Delete(id string) error
	FindByFilter(filter map[string]interface{}) ([]*domain.Course, error)
	SaveReview(review *domain.Review) error
	FindReviewByID(id string) (*domain.Review, error)
	FindReviewsByCourse(courseId string) ([]*domain.Review, error)
}
