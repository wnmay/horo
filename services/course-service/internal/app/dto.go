package app

import (
	"time"

	"github.com/wnmay/horo/services/course-service/internal/domain"
)

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
	CourseID     string
	CustomerID   string
	CustomerName string
	Score        float64
	Title        string
	Description  string
	CreatedAt    time.Time
	DeletedAt    bool
}
