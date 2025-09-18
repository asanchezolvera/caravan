package auth

import (
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

// RegisterUser handles user registration, including password hashing via Argon2.
func (s *Service) RegisterUser(user *models.User) error {

	// Hash user password using Argon2.
	hash := argon2.IDKey([]byte(user.Password), []byte("salt"), 1, 64*1024, 4, 32)

	user.Password = fmt.Sprintf("%x", hash)

	// Create new user in the database.
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

// LoginUser authenticates a user and generates a JWT token upon success.
func (s *Service) LoginUser(user *models.User) (string, error) {

	// Retrieve user from the database by email.
	foundUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil || foundUser == nil {
		return "", fmt.Errorf("Invalid credentials")
	}

	// Verify user password against the stored hash.
	hash := argon2.IDKey([]byte(user.Password), []byte("salt"), 1, 64*1024, 4, 32)
	providedHash := fmt.Sprintf("%x", hash)
	if providedHash != foundUser.Password {
		return "", fmt.Errorf("Invalid credentials")
	}

	// Generate JWT token for the authenticated user.
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

// GetUserProfile retrieves user profile information based on email.
func (s *Service) GetUserProfile(email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve user profile: %w", err)
	}
	return user, nil
}
