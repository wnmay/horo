package db

import (
	"context"

	"github.com/wnmay/horo/services/course-service/internal/domain"

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
