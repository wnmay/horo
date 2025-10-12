package ports

import "github.com/wnmay/horo/services/user-management-service/internal/domain"

type UserRepository interface {
	Save(user domain.User) error
}
