package db

type UserModel struct {
	ID       string `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex"`
	FullName string
	Password string
}