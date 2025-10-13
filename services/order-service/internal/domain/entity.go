package domain

type Person struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}
