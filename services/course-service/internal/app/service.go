package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wnmay/horo/services/course-service/internal/adapters/outbound/db"
	"github.com/wnmay/horo/services/course-service/internal/domain"
	"github.com/wnmay/horo/services/course-service/internal/ports/outbound"
)

type CourseService interface {
	CreateCourse(input CreateCourseInput) (*domain.Course, error)
	GetCourseByID(ctx context.Context, id string) (*domain.Course, error)
	GetCourseByIDWithProphetName(ctx context.Context, id string) (*domain.CourseWithProphetName, error)
	ListCoursesByProphet(ctx context.Context, prophetID string) ([]*domain.CourseWithProphetName, error)
	UpdateCourse(id string, input *domain.UpdateCourseInput) (*domain.Course, error)
	DeleteCourse(id string) error
	FindCoursesByFilter(ctx context.Context, filter CourseFilter) ([]*domain.Course, error)
	CreateReview(input CreateReviewInput) (*domain.Review, error)
	GetReviewByID(id string) (*domain.Review, error)
	ListReviewsByCourse(courseId string) ([]*domain.Review, error)
}

type courseService struct {
	repo          outbound.CourseRepository
	user_provider outbound.UserProvider
}

func (s courseService) GetCourseByID(ctx context.Context, id string) (*domain.Course, error) {
	return s.repo.FindCourseByID(id)
}

func (s courseService) GetCourseByIDWithProphetName(ctx context.Context, id string) (*domain.CourseWithProphetName, error) {
	course, err := s.repo.FindCourseByID(id)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, fmt.Errorf("course not found")
	}

	prophetID := course.ProphetID
	prophetName, err := s.user_provider.GetProphetName(ctx, prophetID)
	if err != nil {
		return nil, err
	}
	courseWithProphetName := &domain.CourseWithProphetName{
		ID:          course.ID,
		ProphetID:   course.ProphetID,
		ProphetName: prophetName,
		CourseName:  course.CourseName,
		CourseType:  course.CourseType,
		Description: course.Description,
		Price:       course.Price,
		Duration:    course.Duration,
		CreatedAt:   course.CreatedAt,
		DeletedAt:   course.DeletedAt,
	}
	return courseWithProphetName, nil
}

func (s courseService) ListCoursesByProphet(ctx context.Context, prophetID string) ([]*domain.CourseWithProphetName, error) {
	prophetName, err := s.user_provider.GetProphetName(ctx, prophetID)
	if err != nil {
		return nil, err
	}
	courses, err := s.repo.FindCoursesByProphet(prophetID)
	if err != nil {
		return nil, err
	}
	courseWithProphetNames := make([]*domain.CourseWithProphetName, 0)
	for _, course := range courses {
		courseWithProphetName := &domain.CourseWithProphetName{
			ID:          course.ID,
			ProphetID:   course.ProphetID,
			ProphetName: prophetName,
			CourseName:  course.CourseName,
			CourseType:  course.CourseType,
			Description: course.Description,
			Price:       course.Price,
			Duration:    course.Duration,
			CreatedAt:   course.CreatedAt,
			DeletedAt:   course.DeletedAt,
		}
		courseWithProphetNames = append(courseWithProphetNames, courseWithProphetName)
	}
	return courseWithProphetNames, nil
}

func NewCourseService(r outbound.CourseRepository, u outbound.UserProvider) CourseService {
	return &courseService{
		repo:          r,
		user_provider: u,
	}
}

func (s *courseService) CreateCourse(input CreateCourseInput) (*domain.Course, error) {
	c := &domain.Course{
		ID:          generateID("COURSE"),
		ProphetID:   input.ProphetID,
		CourseName:  input.CourseName,
		CourseType:  input.CourseType,
		Description: input.Description,
		Price:       input.Price,
		Duration:    input.Duration,
		CreatedAt:   time.Now(),
		DeletedAt:   false,
	}

	if err := s.repo.SaveCourse(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *courseService) UpdateCourse(id string, input *domain.UpdateCourseInput) (*domain.Course, error) {
	updates := make(map[string]interface{})
	if input.CourseName != "" {
		updates["coursename"] = input.CourseName
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Price != nil {
		updates["price"] = input.Price
	}
	if input.Duration != nil {
		updates["duration"] = input.Duration
	}
	return s.repo.Update(id, updates)
}

func (s *courseService) DeleteCourse(id string) error {
	return s.repo.Delete(id)
}

func (s courseService) FindCoursesByFilter(ctx context.Context, filter CourseFilter) ([]*domain.Course, error) {
	prophetNames, err := s.user_provider.GetProphetIDsByNames(ctx, filter.ProphetName)
	if err != nil {
		return nil, err
	}
	prophetIDs := make([]string, 0)
	for _, prophetName := range prophetNames {
		prophetIDs = append(prophetIDs, prophetName.UserID)
	}
	repoFilter := db.CourseFilter{
		CourseName: filter.CourseName,
		ProphetIDs:  prophetIDs,
		Duration:   filter.Duration,
	}
	return s.repo.FindByFilter(ctx, repoFilter)
}

func (s *courseService) CreateReview(input CreateReviewInput) (*domain.Review, error) {
	c := &domain.Review{
		ID:           generateID("REVIEW"),
		CourseId:     input.CourseId,
		CustomerId:   input.CustomerId,
		CustomerName: input.CustomerName,
		Score:        input.Score,
		Title:        input.Title,
		Description:  input.Description,
		CreatedAt:    time.Now(),
		DeletedAt:    false,
	}

	if err := s.repo.SaveReview(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s courseService) GetReviewByID(id string) (*domain.Review, error) {
	return s.repo.FindReviewByID(id)
}

func (s courseService) ListReviewsByCourse(courseId string) ([]*domain.Review, error) {
	return s.repo.FindReviewsByCourse(courseId)
}

func generateID(objType string) string {
	return objType + "-" + uuid.New().String()
}

type CreateCourseInput struct {
	ID          string
	ProphetID   string
	CourseName  string
	CourseType  domain.CourseType
	Description string
	Price       float64
	Duration    domain.DurationEnum
	CreatedAt   time.Time
	DeletedAt   bool
}

type CreateReviewInput struct {
	ID           string
	CourseId     string
	CustomerId   string
	CustomerName string
	Score        float64
	Title        string
	Description  string
	CreatedAt    time.Time
	DeletedAt    bool
}
