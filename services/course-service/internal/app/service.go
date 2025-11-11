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

type courseService struct {
	repo          outbound.CourseRepository
	user_provider outbound.UserProvider
}

func NewCourseService(r outbound.CourseRepository, u outbound.UserProvider) CourseService {
	return &courseService{
		repo:          r,
		user_provider: u,
	}
}

// Simple get by ID
func (s courseService) GetCourseByID(ctx context.Context, id string) (*domain.Course, error) {
	return s.repo.FindCourseByID(ctx, id)
}

// Get full detailed course with prophet name and reviews
func (s courseService) GetCourseDetailByID(ctx context.Context, id string) (*domain.CourseDetail, error) {
	courseDetail, err := s.repo.FindCourseDetailByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if courseDetail == nil {
		return nil, fmt.Errorf("course not found")
	}

	// Enrich with prophet name from user service
	prophetName, err := s.user_provider.GetProphetName(ctx, courseDetail.ProphetID)
	if err != nil {
		return nil, err
	}
	courseDetail.ProphetName = prophetName

	return courseDetail, nil
}

// Create new course
func (s *courseService) CreateCourse(ctx context.Context, input CreateCourseInput) (*domain.Course, error) {
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
		ReviewCount: 0,
		ReviewScore: 0,
	}

	if err := s.repo.SaveCourse(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// Update course
func (s *courseService) UpdateCourse(ctx context.Context, id string, input *domain.UpdateCourseInput) (*domain.Course, error) {
	updates := make(map[string]interface{})
	if input.CourseName != "" {
		updates["coursename"] = input.CourseName
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.Duration != nil {
		updates["duration"] = *input.Duration
	}
	return s.repo.UpdateCourse(ctx, id, updates)
}

// Soft delete
func (s *courseService) DeleteCourse(ctx context.Context, id string) error {
	return s.repo.DeleteCourse(ctx, id)
}

// List all courses for a given prophet, with prophet name attached
func (s courseService) ListCoursesByProphet(ctx context.Context, prophetID string) ([]*domain.CourseWithProphetName, error) {
	prophetName, err := s.user_provider.GetProphetName(ctx, prophetID)
	if err != nil {
		return nil, err
	}

	courses, err := s.repo.FindCoursesByProphet(ctx, prophetID)
	if err != nil {
		return nil, err
	}

	results := make([]*domain.CourseWithProphetName, 0, len(courses))
	for _, course := range courses {
		results = append(results, &domain.CourseWithProphetName{
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
			ReviewCount: course.ReviewCount,
			ReviewScore: course.ReviewScore,
		})
	}
	return results, nil
}

// Find courses by filter (supports filtering and sorting)
func (s courseService) FindCoursesByFilter(ctx context.Context, filter CourseFilter, sort CourseSort) ([]*domain.Course, error) {
	prophets, err := s.user_provider.GetProphetIDsByNames(ctx, filter.ProphetName)
	if err != nil {
		return nil, err
	}

	prophetIDs := make([]string, len(prophets))
	for i, p := range prophets {
		prophetIDs[i] = p.UserID
	}

	repoFilter := db.CourseFilter{
		CourseName: filter.CourseName,
		ProphetIDs: prophetIDs,
		Duration:   filter.Duration,
		CourseType: string(filter.CourseType),
	}

	repoSort := db.CourseSort{
		SortBy: string(sort.SortBy),
		Order:  sort.Order,
	}

	return s.repo.FindByFilter(ctx, repoFilter, repoSort)
}

// Create a new review and automatically update course’s denormalized score
func (s *courseService) CreateReview(ctx context.Context, input CreateReviewInput) (*domain.Review, error) {
	review := &domain.Review{
		ID:           generateID("REVIEW"),
		CourseID:     input.CourseID,
		CustomerID:   input.CustomerID,
		CustomerName: input.CustomerName,
		Score:        input.Score,
		Title:        input.Title,
		Description:  input.Description,
		CreatedAt:    time.Now(),
		DeletedAt:    false,
	}

	if err := s.repo.SaveReview(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

// Get a single review by ID
func (s courseService) GetReviewByID(ctx context.Context, id string) (*domain.Review, error) {
	return s.repo.FindReviewByID(ctx, id)
}

// List all reviews for a given course
func (s courseService) ListReviewsByCourse(ctx context.Context, courseID string) ([]*domain.Review, error) {
	return s.repo.FindReviewsByCourse(ctx, courseID)
}

// helper
func generateID(prefix string) string {
	return prefix + "-" + uuid.New().String()
}
