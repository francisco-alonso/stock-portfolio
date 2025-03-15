package services

import (
	"github.com/francisco-alonso/go-template/internal/domain"
	"github.com/francisco-alonso/go-template/internal/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, email string) error {
	user := domain.User{
		Username: username,
		Email:    email,
	}
	
	return s.repo.Create(user)
}