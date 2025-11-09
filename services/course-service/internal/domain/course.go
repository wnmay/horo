// internal/domain/course.go
package domain

import "time"

type Course struct {
	ID          string       `bson:"id"`
	ProphetID   string       `bson:"prophet_id"`
	ProphetName string       `bson:"prophetname"`
	CourseName  string       `bson:"coursename"`
	Description string       `bson:"description"`
	Price       float64      `bson:"price"`
	Duration    DurationEnum `bson:"duration"`
	CreatedAt   time.Time    `bson:"created_time"`
	DeletedAt   bool         `bson:"deleted_at"`
}

type UpdateCourseInput struct {
	CourseName  string        `json:"coursename,omitempty"`
	Description string        `json:"description,omitempty"`
	Price       *float64      `json:"price,omitempty"`
	Duration    *DurationEnum `json:"duration,omitempty"`
	DeletedAt   bool          `json:"deleted_at,omitempty"`
}

type DurationEnum int32

const (
	DurationEnum_15 DurationEnum = 15
	DurationEnum_30 DurationEnum = 30
	DurationEnum_60 DurationEnum = 60
)
