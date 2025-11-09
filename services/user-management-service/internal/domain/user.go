// internal/domain/user.go
package domain

type User struct {
	ID       string
	FullName string
	Email    string
	Role     string
}

type Prophet struct {
	User    *User
	Balance float64
}
