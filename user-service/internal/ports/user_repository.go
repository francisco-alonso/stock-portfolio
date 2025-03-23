package ports

import "github.com/francisco-alonso/stock-portfolio/user-service/internal/domain"

type UserRepository interface {
	Create(user domain.User) error
}