package ports

import "github.com/francisco-alonso/go-template/internal/domain"

type UserRepository interface {
	Create(user domain.User) error
}