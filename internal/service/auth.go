package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	todo "github.com/slemeshaev/todo/internal/models"
	"github.com/slemeshaev/todo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 12 * time.Hour

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUsernameTaken      = errors.New("username already taken")
	dummyHash             = mustHash("")
)

type AuthService struct {
	repo       repository.Authorization
	signingKey []byte
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization, signingKey string) *AuthService {
	return &AuthService{repo: repo, signingKey: []byte(signingKey)}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}
	user.Password = string(hash)

	id, err := s.repo.CreateUser(user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return 0, ErrUsernameTaken
		}
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// Spend the same time as a real check so the response time does not reveal unknown usernames.
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		user.Id,
	})

	return token.SignedString(s.signingKey)
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{},
		func(token *jwt.Token) (any, error) {
			return s.signingKey, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func mustHash(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("generate dummy hash: %v", err))
	}

	return hash
}
