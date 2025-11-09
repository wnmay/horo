package db

import (
	"context"
	"strconv"

	"github.com/wnmay/horo/services/course-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCourseRepo struct {
	col *mongo.Collection
}

func NewMongoCourseRepo(db *mongo.Database) *MongoCourseRepo {
	return &MongoCourseRepo{col: db.Collection("courses")}
}

func (r *MongoCourseRepo) SaveCourse(course *domain.Course) error {
	_, err := r.col.InsertOne(context.TODO(), course)
	return err
}

func (r *MongoCourseRepo) FindCourseByID(id string) (*domain.Course, error) {
	var c domain.Course
	err := r.col.FindOne(context.TODO(), bson.M{"id": id, "deleted_at": false}).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MongoCourseRepo) FindCoursesByProphet(prophetID string) ([]*domain.Course, error) {
	cur, err := r.col.Find(context.TODO(), bson.M{"prophet_id": prophetID, "deleted_at": false})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var courses []*domain.Course
	for cur.Next(context.TODO()) {
		var c domain.Course
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}
	return courses, nil
}

func (r *MongoCourseRepo) Update(id string, updates map[string]interface{}) (*domain.Course, error) {
	_, err := r.col.UpdateOne(context.TODO(), bson.M{"id": id, "deleted_at": false}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return r.FindCourseByID(id)
}

func (r *MongoCourseRepo) Delete(id string) error {
	update := bson.M{"$set": bson.M{"deleted_at": true}}
	_, err := r.col.UpdateOne(context.TODO(), bson.M{"id": id}, update)
	return err
}

func (r *MongoCourseRepo) FindByFilter(filter map[string]interface{}) ([]*domain.Course, error) {
	query := bson.M{
		"deleted_at": false,
		"id":         bson.M{"$regex": "^COURSE-"},
	}
	for key, val := range filter {
		switch key {
		case "coursename":
			query["coursename"] = bson.M{"$regex": val, "$options": "i"}
		case "prophetname":
			query["prophetname"] = bson.M{"$regex": val, "$options": "i"}
		case "duration":
			if valStr, ok := val.(string); ok {
				if i, err := strconv.Atoi(valStr); err == nil {
					query["duration"] = i
				}
			} else {
				query["duration"] = val
			}
		}
	}

	cursor, err := r.col.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var courses []*domain.Course
	if err := cursor.All(context.TODO(), &courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *MongoCourseRepo) SaveReview(review *domain.Review) error {
	_, err := r.col.InsertOne(context.TODO(), review)
	return err
}

func (r *MongoCourseRepo) FindReviewByID(id string) (*domain.Review, error) {
	var rv domain.Review
	err := r.col.FindOne(context.TODO(), bson.M{"id": id, "deleted_at": false}).Decode(&rv)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

func (r *MongoCourseRepo) FindReviewsByCourse(courseId string) ([]*domain.Review, error) {
	cur, err := r.col.Find(context.TODO(), bson.M{"course_id": courseId, "deleted_at": false})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var reviews []*domain.Review
	for cur.Next(context.TODO()) {
		var rv domain.Review
		if err := cur.Decode(&rv); err != nil {
			return nil, err
		}
		reviews = append(reviews, &rv)
	}
	return reviews, nil
}
