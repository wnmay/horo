package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type CourseFilter struct {
	CourseName string
	ProphetIDs []string
	Duration   string
	CourseType string
}

type CourseSort struct {
	SortBy string // "price", "review_score"
	Order  string // "asc", "desc"
}

func BuildMongoQuery(filter CourseFilter, sort CourseSort) (bson.M, bson.D, error) {
	mongoFilter := bson.M{}

	if filter.CourseName != "" {
		mongoFilter["coursename"] = bson.M{
			"$regex":   strings.TrimSpace(filter.CourseName),
			"$options": "i",
		}
	}

	if len(filter.ProphetIDs) > 0 {
		mongoFilter["prophet_id"] = bson.M{"$in": filter.ProphetIDs}
	}

	if filter.Duration != "" {
		durationInt, err := strconv.Atoi(filter.Duration)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid duration: %v", err)
		}
		mongoFilter["duration"] = durationInt
	}

	if filter.CourseType != "" {
		mongoFilter["coursetype"] = bson.M{
			"$regex":   strings.TrimSpace(filter.CourseType),
			"$options": "i",
		}
	}

	mongoSort := bson.D{}
	if sort.SortBy != "" {
		order := 1
		if strings.ToLower(sort.Order) == "desc" {
			order = -1
		}
		mongoSort = append(mongoSort, bson.E{Key: sort.SortBy, Value: order})
	}

	log.Printf("Mongo Filter: %+v\n", mongoFilter)
	log.Printf("Mongo Sort: %+v\n", mongoSort)

	return mongoFilter, mongoSort, nil
}
