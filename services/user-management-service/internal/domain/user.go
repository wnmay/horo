// internal/domain/user.go
package domain

type User struct {
	ID       string
	FullName string
	Email    string
	Password string // hashed
	Role     string
}


