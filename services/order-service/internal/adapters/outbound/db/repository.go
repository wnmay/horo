package db

import (
	"context"

	"github.com/wnmay/horo/services/order-service/internal/domain"
	"github.com/wnmay/horo/services/order-service/internal/ports/outbound"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type personDoc struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

type MongoPersonRepository struct {
	col *mongo.Collection
}

var _ outbound.PersonRepository = (*MongoPersonRepository)(nil)

func NewMongoPersonRepository(db *mongo.Database) *MongoPersonRepository {
	col := db.Collection("order")
	_, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetName("idx_name_asc"),
	})
	return &MongoPersonRepository{col: col}
}

func (r *MongoPersonRepository) Save(p domain.Person) error {
	_, err := r.col.InsertOne(context.Background(), personDoc{ID: p.ID, Name: p.Name})
	return err
}

func (r *MongoPersonRepository) GetAll() ([]domain.Person, error) {
	cur, err := r.col.Find(context.Background(), bson.D{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var out []domain.Person
	for cur.Next(context.Background()) {
		var d personDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		out = append(out, domain.Person{ID: d.ID, Name: d.Name})
	}
	return out, cur.Err()
}
