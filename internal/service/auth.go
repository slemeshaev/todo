package service

import (
	todo "github.com/slemeshaev/todo/internal/models"
	"github.com/slemeshaev/todo/internal/repository"
)

type AuthService struct {
	repo repository.Repository
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	return s.repo.CreateUser(user)
}
