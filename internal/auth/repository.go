package auth

import (
	"fmt"

	"gorm.io/gorm"

	"caravan/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser creates a new user in the database.
func (r *Repository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

// GetUserByEmail retrieves a user by email from the database.
func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	result := r.db.Where("email = ?", email).First(user)

	// If user not found, return nil User and nil Error.
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", result.Error)
	}
	return user, nil
}
