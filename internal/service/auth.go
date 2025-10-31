package service

import (
	"crypto/sha1"
	"fmt"

	todo "github.com/slemeshaev/todo/internal/models"
	"github.com/slemeshaev/todo/internal/repository"
)

const salt = "jdkafjajl84uo8afsf44fafadfg"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
