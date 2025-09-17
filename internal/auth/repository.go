package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"caravan/internal/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateUser creates a new user in the database.
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by email from the database.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx, "SELECT email, password_hash FROM users WHERE email = $1", email).Scan(&user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}
