package db

import (
	"context"
	"log"

	"github.com/wnmay/horo/services/course-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCourseRepo struct {
	courseCol *mongo.Collection
	reviewCol *mongo.Collection
}

func NewMongoCourseRepo(db *mongo.Database) *MongoCourseRepo {
	return &MongoCourseRepo{
		courseCol: db.Collection("courses"),
		reviewCol: db.Collection("reviews"),
	}
}

func (r *MongoCourseRepo) SaveCourse(ctx context.Context, course *domain.Course) error {
	_, err := r.courseCol.InsertOne(ctx, course)
	return err
}

func (r *MongoCourseRepo) FindCourseByID(ctx context.Context, id string) (*domain.Course, error) {
	var c domain.Course
	err := r.courseCol.FindOne(ctx, bson.M{"id": id, "deleted_at": false}).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MongoCourseRepo) FindCourseDetailByID(ctx context.Context, id string) (*domain.CourseDetail, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"id": id, "deleted_at": false}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "reviews",
			"localField":   "id",
			"foreignField": "course_id",
			"as":           "reviews",
		}}},
		// Optional: ensure fields exist, but don’t recalc
		{{Key: "$addFields", Value: bson.M{
			"review_count": bson.M{"$ifNull": []interface{}{"$review_count", 0}},
			"review_score": bson.M{"$ifNull": []interface{}{"$review_score", 0.0}},
		}}},
	}

	cur, err := r.courseCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return nil, mongo.ErrNoDocuments
	}

	var detail domain.CourseDetail
	if err := cur.Decode(&detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

func (r *MongoCourseRepo) FindCoursesByProphet(ctx context.Context, prophetID string) ([]*domain.Course, error) {
	cur, err := r.courseCol.Find(ctx, bson.M{"prophet_id": prophetID, "deleted_at": false})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var courses []*domain.Course
	for cur.Next(ctx) {
		var c domain.Course
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}
	return courses, nil
}

func (r *MongoCourseRepo) UpdateCourse(ctx context.Context, id string, updates map[string]interface{}) (*domain.Course, error) {
	_, err := r.courseCol.UpdateOne(ctx, bson.M{"id": id, "deleted_at": false}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return r.FindCourseByID(ctx, id)
}

func (r *MongoCourseRepo) DeleteCourse(ctx context.Context, id string) error {
	update := bson.M{"$set": bson.M{"deleted_at": true}}
	_, err := r.courseCol.UpdateOne(ctx, bson.M{"id": id}, update)
	return err
}

func (r *MongoCourseRepo) FindByFilter(ctx context.Context, filter CourseFilter, sort CourseSort) ([]*domain.Course, error) {
	filterMongo, sortMongo, err := BuildMongoQuery(filter, sort)
	if err != nil {
		log.Printf("Error building MongoDB query: %v", err)
		return nil, err
	}

	// If no filter criteria are provided, return all active courses
	if filter.SearchTerm == "" && len(filter.ProphetIDs) == 0 && filter.Duration == "" && filter.CourseType == "" {
		filterMongo = bson.M{"deleted_at": false}
		log.Println("No filters provided — returning all courses.")
	}

	opts := options.Find()
	if len(sortMongo) > 0 {
		opts.SetSort(sortMongo)
	}

	cursor, err := r.courseCol.Find(ctx, filterMongo, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var courses []*domain.Course
	for cursor.Next(ctx) {
		var c domain.Course
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}

	return courses, nil
}

func (r *MongoCourseRepo) SaveReview(ctx context.Context, review *domain.Review) error {
	// Insert review
	if _, err := r.reviewCol.InsertOne(ctx, review); err != nil {
		return err
	}

	// Recalculate review count + average
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"course_id": review.CourseID, "deleted_at": false}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$course_id",
			"count": bson.M{"$sum": 1},
			"avg":   bson.M{"$avg": "$score"},
		}}},
	}

	cur, err := r.reviewCol.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	var agg struct {
		ID    string  `bson:"_id"`
		Count int     `bson:"count"`
		Avg   float64 `bson:"avg"`
	}

	if cur.Next(ctx) {
		if err := cur.Decode(&agg); err != nil {
			return err
		}

		update := bson.M{"$set": bson.M{
			"review_count": agg.Count,
			"review_score": agg.Avg,
		}}
		_, err = r.courseCol.UpdateOne(ctx, bson.M{"id": review.CourseID}, update)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MongoCourseRepo) FindReviewsByCourse(ctx context.Context, courseID string) ([]*domain.Review, error) {
	cur, err := r.reviewCol.Find(ctx, bson.M{"course_id": courseID, "deleted_at": false})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var reviews []*domain.Review
	for cur.Next(ctx) {
		var rv domain.Review
		if err := cur.Decode(&rv); err != nil {
			return nil, err
		}
		reviews = append(reviews, &rv)
	}
	return reviews, nil
}

func (r *MongoCourseRepo) FindReviewByID(ctx context.Context, id string) (*domain.Review, error) {
	var rv domain.Review
	err := r.reviewCol.FindOne(ctx, bson.M{"id": id, "deleted_at": false}).Decode(&rv)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

func (r *MongoCourseRepo) FindPopularCourses(ctx context.Context, limit int) ([]*domain.Course, error) {
	cur, err := r.courseCol.Find(ctx, bson.M{"deleted_at": false}, options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "review_score", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var courses []*domain.Course
	for cur.Next(ctx) {
		var c domain.Course
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}
	return courses, nil
}
