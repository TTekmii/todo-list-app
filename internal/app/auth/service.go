package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
)

type Service struct {
	repo       repo.Authorization
	jwtSecret  string
	tokenTTL   time.Duration
	bcryptCost int
}

func NewService(authRepo repo.Authorization, jwtSecret string, tokenTTL time.Duration, bcryptCost int) *Service {
	return &Service{
		repo:       authRepo,
		jwtSecret:  jwtSecret,
		tokenTTL:   tokenTTL,
		bcryptCost: bcryptCost,
	}
}

func (s *Service) Register(ctx context.Context, username, name, password string) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password:%w", err)
	}

	user := model.User{
		Name:         name,
		Username:     username,
		PasswordHash: string(hash),
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.generateToken(user.ID, user.Username)
}

func (s *Service) generateToken(userID int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"name": username,
		"exp":  time.Now().Add(s.tokenTTL).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}
