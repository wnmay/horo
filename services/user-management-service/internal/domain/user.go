// internal/domain/user.go
package domain

type User struct {
	ID       string
	FullName string
	Email    string
	Role     string
}

type ProphetName struct {
	UserID   string
	ProphetName string
}