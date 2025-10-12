// internal/domain/user.go
package domain

type User struct {
	ID       string
	FullName string
	Email    string
	Password string // hashed
	Role     string
}

type Claims struct {
	UserID string
	Email  string
	Role   string
}
