package app

import (
	"context"

	"github.com/wnmay/horo/services/course-service/internal/domain"
)

type CourseService interface {
	CreateCourse(ctx context.Context, input CreateCourseInput) (*domain.Course, error)
	GetCourseByID(ctx context.Context, id string) (*domain.Course, error)
	GetCourseDetailByID(ctx context.Context, id string) (*domain.CourseDetail, error)
	ListCoursesByProphet(ctx context.Context, prophetID string) ([]*domain.CourseWithProphetName, error)
	UpdateCourse(ctx context.Context, id string, input *domain.UpdateCourseInput) (*domain.Course, error)
	DeleteCourse(ctx context.Context, id string) error
	FindCoursesByFilter(ctx context.Context, filter CourseFilter, sort CourseSort) ([]*domain.CourseWithProphetName, error)
	CreateReview(ctx context.Context, input CreateReviewInput) (*domain.Review, error)
	GetReviewByID(ctx context.Context, id string) (*domain.Review, error)
	ListReviewsByCourse(ctx context.Context, courseId string) ([]*domain.Review, error)
	ListPopularCourses(ctx context.Context, limit int) ([]*domain.CourseWithProphetName, error)
}
