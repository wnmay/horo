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
	UpdateCourse(id string, input *domain.UpdateCourseInput) (*domain.Course, error)
	DeleteCourse(id string) error
	FindCoursesByFilter(filter map[string]interface{}) ([]*domain.Course, error)
	CreateReview(input CreateReviewInput) (*domain.Review, error)
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
		ID:          generateID("COURSE"),
		ProphetID:   input.ProphetID,
		ProphetName: input.ProphetName,
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

func (s *courseService) FindCoursesByFilter(filter map[string]interface{}) ([]*domain.Course, error) {
	return s.repo.FindByFilter(filter)
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

func generateID(objType string) string {
	return objType + "-" + uuid.New().String()
}

type CreateCourseInput struct {
	ID          string
	ProphetID   string
	ProphetName string
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
