// internal/domain/user.go
package domain

type UserRole string

const (
	USER_ROLE_PROPHET  UserRole = "prophet"
	USER_ROLE_CUSTOMER UserRole = "customer"
	USER_ROLE_UNKNOWN  UserRole = "unknown"
)

type User struct {
	ID       string
	FullName string
	Email    string
	Role     string
}

type ProphetName struct {
	UserID      string
	ProphetName string
}

type UserName struct {
	UserID   string
	UserName string
	UserRole UserRole
}
