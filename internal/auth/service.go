package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/argon2"

	"caravan/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) RegisterUser(ctx context.Context, user *models.User) error {
	hash := argon2.IDKey([]byte(user.Password), []byte("salt"), 1, 64*1024, 4, 32)
	user.Password = fmt.Sprintf("%x", hash)

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) LoginUser(ctx context.Context, user *models.User) (string, error) {
	foundUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", fmt.Errorf("Invalid credentials")
	}

	hash := argon2.IDKey([]byte(user.Password), []byte("salt"), 1, 64*1024, 4, 32)
	providedHash := fmt.Sprintf("%x", hash)
	if providedHash != foundUser.Password {
		return "", fmt.Errorf("Invalid credentials")
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &models.Claims{
		Email: foundUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token")
	}

	return tokenString, nil
}
