// internal/domain/course.go
package domain

import (
	"time"
)

type Course struct {
	ID          string       `bson:"id" json:"id"`
	ProphetID   string       `bson:"prophet_id" json:"prophet_id"`
	CourseName  string       `bson:"coursename" json:"coursename"`
	CourseType  CourseType   `bson:"coursetype" json:"coursetype"`
	Description string       `bson:"description" json:"description"`
	Price       float64      `bson:"price" json:"price"`
	Duration    DurationEnum `bson:"duration" json:"duration"`
	CreatedAt   time.Time    `bson:"created_time" json:"created_time"`
	DeletedAt   bool         `bson:"deleted_at" json:"deleted_at"`
}

type CourseWithProphetName struct {
	ID          string       `bson:"id" json:"id"`
	ProphetID   string       `bson:"prophet_id" json:"prophet_id"`
	ProphetName string       `bson:"prophetname" json:"prophetname"`
	CourseName  string       `bson:"coursename" json:"coursename"`
	CourseType  CourseType   `bson:"coursetype" json:"coursetype"`
	Description string       `bson:"description" json:"description"`
	Price       float64      `bson:"price" json:"price"`
	Duration    DurationEnum `bson:"duration" json:"duration"`
	CreatedAt   time.Time    `bson:"created_time" json:"created_time"`
	DeletedAt   bool         `bson:"deleted_at" json:"deleted_at"`
}

type Review struct {
	ID           string    `bson:"id"            json:"id"`
	CourseId     string    `bson:"course_id"     json:"course_id"`
	CustomerId   string    `bson:"customer_id"   json:"customer_id"`
	CustomerName string    `bson:"customername"  json:"customername"`
	Score        float64   `bson:"score"         json:"score"`
	Title        string    `bson:"title"         json:"title"`
	Description  string    `bson:"description"   json:"description"`
	CreatedAt    time.Time `bson:"created_at"    json:"created_at"`
	DeletedAt    bool      `bson:"deleted_at"    json:"deleted_at"`
}

type UpdateCourseInput struct {
	CourseName  string        `bson:"coursename,omitempty"  json:"coursename,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	Price       *float64      `bson:"price,omitempty"       json:"price,omitempty"`
	Duration    *DurationEnum `bson:"duration,omitempty"    json:"duration,omitempty"`
	DeletedAt   bool          `bson:"deleted_at,omitempty"  json:"deleted_at,omitempty"`
}

type UpdateReviewInput struct {
	Score       float64 `bson:"score,omitempty"       json:"score,omitempty"`
	Title       string  `bson:"title,omitempty"       json:"title,omitempty"`
	Description string  `bson:"description,omitempty" json:"description,omitempty"`
	DeletedAt   bool    `bson:"deleted_at,omitempty"  json:"deleted_at,omitempty"`
}

type DurationEnum int32

const (
	DurationEnum_15 DurationEnum = 15
	DurationEnum_30 DurationEnum = 30
	DurationEnum_60 DurationEnum = 60
)

type CourseType string

const (
	CourseType_Love  CourseType = "Love"
	CourseType_Work  CourseType = "Work"
	CourseType_Money CourseType = "Money"
	CourseType_Luck  CourseType = "Luck"
	CourseType_Study CourseType = "Study"
)
