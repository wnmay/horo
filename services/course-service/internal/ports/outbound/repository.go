package outbound

import (
	"context"

	"github.com/wnmay/horo/services/course-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/course-service/internal/domain"
)

type CourseRepository interface {
	//Course
	SaveCourse(ctx context.Context, course *domain.Course) error
	FindCourseByID(ctx context.Context, id string) (*domain.Course, error)
	FindCoursesByProphet(ctx context.Context, prophetID string) ([]*domain.Course, error)
	UpdateCourse(ctx context.Context, id string, updates map[string]interface{}) (*domain.Course, error)
	DeleteCourse(ctx context.Context, id string) error
	
	//Filter, sort
	FindByFilter(ctx context.Context, filter db.CourseFilter, sort db.CourseSort) ([]*domain.Course, error)

	//Review
	SaveReview(ctx context.Context, review *domain.Review) error
	FindReviewByID(ctx context.Context, id string) (*domain.Review, error)
	FindReviewsByCourse(ctx context.Context, courseID string) ([]*domain.Review, error)

	//Course with review
	FindCourseDetailByID(ctx context.Context, id string) (*domain.CourseDetail, error)
}
