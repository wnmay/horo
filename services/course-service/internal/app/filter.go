package app

import (
	"strings"

	"github.com/wnmay/horo/services/course-service/internal/domain"
)

type CourseFilter struct {
	CourseName  string
	ProphetName string
	Duration    string
	CourseType  domain.CourseType
}

type CourseSort struct {
	SortBy SortType
	Order  string
}

type SortType string

const (
	SortBy_Price  SortType = "price"
	SortBy_Score SortType = "score"
)

// Map plain string to domain.CourseType
func ParseCourseType(input string) domain.CourseType {
	switch strings.ToLower(input) {
	case string(domain.CourseType_Love):
		return domain.CourseType_Love
	case string(domain.CourseType_Work):
		return domain.CourseType_Work
	case string(domain.CourseType_Money):
		return domain.CourseType_Money
	case string(domain.CourseType_Luck):
		return domain.CourseType_Luck
	case string(domain.CourseType_Study):
		return domain.CourseType_Study
	case string(domain.CourseType_Growth):
		return domain.CourseType_Growth
	default:
		return "" // invalid, treated as no filter
	}
}

// Map plain string to SortType
func ParseSortType(input string) SortType {
	switch strings.ToLower(input) {
	case string(SortBy_Price):
		return SortBy_Price
	case string(SortBy_Score):
		return SortBy_Score
	default:
		return ""
	}
}
