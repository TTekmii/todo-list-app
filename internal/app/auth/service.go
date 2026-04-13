package auth

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
	"github.com/golang-jwt/jwt"
)

// TODO: Replace sha1 with bcrypt
// TODO: Move signingKey, salt, and tokenTTL to the config
// TODO: Add context.Context to all methods
// TODO: Correct the JWT-claims structure

const (
	salt       = "fhuysdf85139sdgjkfsd"
	signingKey = "sfduu#5f#41sdf#562sd"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repo.Authorization
}

func NewAuthService(repo repo.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

// func (s *Service) Register(ctx context.Context, input dto.RegisterInput) (model.PublicUser, error) {
//     // 1. Валидация (если не используете binding-теги)
//     // 2. Хеширование пароля
//     hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
//     if err != nil {
//         return model.PublicUser{}, err
//     }

//     // 3. Создание доменной модели
//     user := model.User{
//         Name:         input.Name,
//         Username:     input.Username,
//         PasswordHash: string(hash),
//     }

//     // 4. Сохранение через репозиторий
//     id, err := s.repo.CreateUser(ctx, user)
//     if err != nil {
//         return model.PublicUser{}, err
//     }

//     // 5. Возврат публичной версии
//     return model.PublicUser{
//         ID:       id,
//         Name:     user.Name,
//         Username: user.Username,
//     }, nil
// }
