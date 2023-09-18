package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository"
)

type AuthRepository struct {
	pool Driver
}

func NewAuthRepository(pool Driver) *AuthRepository {
	return &AuthRepository{pool}
}

func (r *AuthRepository) CheckUserExists(ctx context.Context, id int16) error {
	var name string
	err := r.pool.QueryRow(ctx, `
		SELECT name 
		FROM users
		WHERE id=$1 AND disabled_at IS NULL`, id).Scan(&name)
	if err != nil {
		return fmt.Errorf("selecting user: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, encoded_password, salt
		FROM users
		WHERE email=$1 OR login=$1 AND disabled_at IS NULL`, login).Scan(
		&user.ID,
		&user.EncodedPassword,
		&user.Salt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	return &user, nil
}
