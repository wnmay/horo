package domain

import "time"

//
// ─── COURSE STRUCTS ────────────────────────────────────────────────────────────
//

// The main course entity (stored in "courses" collection)
type Course struct {
	ID           string       `bson:"id"            json:"id"`
	ProphetID    string       `bson:"prophet_id"    json:"prophet_id"`
	CourseName   string       `bson:"coursename"    json:"coursename"`
	CourseType   CourseType   `bson:"coursetype"    json:"coursetype"`
	Description  string       `bson:"description"   json:"description"`
	Price        float64      `bson:"price"         json:"price"`
	Duration     DurationEnum `bson:"duration"      json:"duration"`
	CreatedAt    time.Time    `bson:"created_time"  json:"created_time"`
	DeletedAt    bool         `bson:"deleted_at"    json:"deleted_at"`

	// Denormalized fields — automatically updated when reviews change
	ReviewCount int     `bson:"review_count" json:"review_count"`
	ReviewScore float64 `bson:"review_score" json:"review_score"`
}

// Course with Prophet Name (useful for listing / joins)
type CourseWithProphetName struct {
	ID           string       `bson:"id"            json:"id"`
	ProphetID    string       `bson:"prophet_id"    json:"prophet_id"`
	ProphetName  string       `bson:"prophetname"   json:"prophetname"`
	CourseName   string       `bson:"coursename"    json:"coursename"`
	CourseType   CourseType   `bson:"coursetype"    json:"coursetype"`
	Description  string       `bson:"description"   json:"description"`
	Price        float64      `bson:"price"         json:"price"`
	Duration     DurationEnum `bson:"duration"      json:"duration"`
	CreatedAt    time.Time    `bson:"created_time"  json:"created_time"`
	DeletedAt    bool         `bson:"deleted_at"    json:"deleted_at"`
	ReviewCount  int          `bson:"review_count"  json:"review_count"`
	ReviewScore  float64      `bson:"review_score"  json:"review_score"`
}

// CourseDetail — returned when joining course + reviews
type CourseDetail struct {
	ID           string       `bson:"id"            json:"id"`
	ProphetID    string       `bson:"prophet_id"    json:"prophet_id"`
	ProphetName  string       `bson:"prophetname"   json:"prophetname"`
	CourseName   string       `bson:"coursename"    json:"coursename"`
	CourseType   CourseType   `bson:"coursetype"    json:"coursetype"`
	Description  string       `bson:"description"   json:"description"`
	Price        float64      `bson:"price"         json:"price"`
	Duration     DurationEnum `bson:"duration"      json:"duration"`
	CreatedAt    time.Time    `bson:"created_time"  json:"created_time"`
	DeletedAt    bool         `bson:"deleted_at"    json:"deleted_at"`

	// Aggregated fields
	ReviewCount int       `bson:"review_count" json:"review_count"`
	ReviewScore float64   `bson:"review_score" json:"review_score"`

	// Embedded array of reviews (from $lookup)
	Reviews []*Review `bson:"reviews" json:"reviews"`
}

//
// ─── REVIEW STRUCTS ────────────────────────────────────────────────────────────
//

// Stored in "reviews" collection
type Review struct {
	ID           string    `bson:"id"            json:"id"`
	CourseID     string    `bson:"course_id"     json:"course_id"`
	CustomerID   string    `bson:"customer_id"   json:"customer_id"`
	CustomerName string    `bson:"customername"  json:"customername"`
	Score        float64   `bson:"score"         json:"score"`
	Title        string    `bson:"title"         json:"title"`
	Description  string    `bson:"description"   json:"description"`
	CreatedAt    time.Time `bson:"created_at"    json:"created_at"`
	DeletedAt    bool      `bson:"deleted_at"    json:"deleted_at"`
}

//
// ─── UPDATE INPUT STRUCTS ─────────────────────────────────────────────────────
//

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

//
// ─── ENUMS ────────────────────────────────────────────────────────────────────
//

type DurationEnum int32

const (
	DurationEnum_15 DurationEnum = 15
	DurationEnum_30 DurationEnum = 30
	DurationEnum_60 DurationEnum = 60
)

type CourseType string

const (
	CourseType_Love   CourseType = "love"
	CourseType_Work   CourseType = "work"
	CourseType_Money  CourseType = "money"
	CourseType_Luck   CourseType = "luck"
	CourseType_Study  CourseType = "study"
	CourseType_Growth CourseType = "personal_growth"
)
