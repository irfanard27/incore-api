package pgsql

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (name, email, password, created_at)
		VALUES (:name, :email, :password, NOW())
	`
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE email = :email
	`
	params := map[string]interface{}{
		"email": email,
	}
	var user entity.User
	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	err = rows.StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return &user, nil
}
