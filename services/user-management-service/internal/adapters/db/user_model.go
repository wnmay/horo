package db

// Model for saving into gorm
type UserModel struct {
	UserID   string `bson:"user_id,omitempty" gorm:"primaryKey"`
	FullName string `bson:"fullname"`
	Email    string `bson:"email"`
	Role     string `bson:"role"` // "prophet" or "customer"
}

type ProphetModel struct {
	User    UserModel `bson:",inline"`
	Balance float64   `bson:"balance"`
}
