package db

import (
	"context"

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

func (r *MongoCourseRepo) Save(course *domain.Course) error {
	_, err := r.col.InsertOne(context.TODO(), course)
	return err
}

func (r *MongoCourseRepo) FindByID(id string) (*domain.Course, error) {
	var c domain.Course
	err := r.col.FindOne(context.TODO(), bson.M{"id": id}).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MongoCourseRepo) FindAllByProphet(prophetID string) ([]*domain.Course, error) {
	cur, err := r.col.Find(context.TODO(), bson.M{"prophet_id": prophetID})
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
