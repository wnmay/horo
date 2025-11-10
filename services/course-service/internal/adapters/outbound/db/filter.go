package db

import (
	"log"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type CourseFilter struct {
	CourseName string   `bson:"coursename"`
	ProphetIDs []string `bson:"prophet_ids"` // filter input: array of prophet IDs
	Duration   string   `bson:"duration"`
}

func (f CourseFilter) BuildMongoFilter() (bson.M, error) {
	filter := bson.M{}

	// Case-insensitive partial match for course name
	if f.CourseName != "" {
		filter["coursename"] = bson.M{
			"$regex":   strings.TrimSpace(f.CourseName),
			"$options": "i",
		}
	}

	// Prophet ID: field in Mongo is prophet_id, filter values are ProphetIDs
	log.Println("ProphetIDs", f.ProphetIDs)
	
	if len(f.ProphetIDs) > 0 {
		filter["prophet_id"] = bson.M{"$in": f.ProphetIDs}
	}

	// Duration: numeric match
	if f.Duration != "" {
		durationInt, err := strconv.Atoi(f.Duration)
		if err != nil {
			return nil, err
		}
		filter["duration"] = durationInt
	}

	return filter, nil
}
