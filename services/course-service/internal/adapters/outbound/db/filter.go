package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type CourseFilter struct {
	SearchTerm string
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

	// Always filter out deleted courses
	mongoFilter["deleted_at"] = false

	// Handle SearchTerm with OR logic: match course name OR prophet ID
	if filter.SearchTerm != "" && len(filter.ProphetIDs) > 0 {
		// If we have both search term and prophet IDs, use OR logic
		log.Println("Building OR query with SearchTerm and ProphetIDs")
		mongoFilter["$or"] = []bson.M{
			{
				"coursename": bson.M{
					"$regex":   strings.TrimSpace(filter.SearchTerm),
					"$options": "i",
				},
			},
			{
				"prophet_id": bson.M{"$in": filter.ProphetIDs},
			},
		}
	} else if filter.SearchTerm != "" {
		// Only search term, search by course name
		log.Println("Building query with SearchTerm")
		mongoFilter["coursename"] = bson.M{
			"$regex":   strings.TrimSpace(filter.SearchTerm),
			"$options": "i",
		}
	} else if len(filter.ProphetIDs) > 0 {
		log.Println("Building query with ProphetIDs")
		// Only prophet IDs, search by prophet ID
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

	return mongoFilter, mongoSort, nil
}
