package postgresql

import (
	"context"
	"errors"
	"fmt"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	pool Driver
}

func NewUserRepository(pool Driver) *UserRepository {
	return &UserRepository{pool}
}

func (r *UserRepository) GetUser(ctx context.Context, id int16) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, name, surname, login, created_at, role, position
                                 FROM users
                                 WHERE id = $1`, id)

	var user domain.User
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Login,
		&user.CreatedAt,
		&user.Role,
		&user.Position,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}

	user.RoleName = user.Role.String()
	user.PositionName = user.Position.String()

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (int16, error) {
	var userId int16

	err := r.pool.QueryRow(ctx,
		`INSERT INTO users(name, surname, login, password, role, position)
             VALUES($1, $2, $3, $4, $5, $6)
             RETURNING  id`,
		user.Name,
		user.Surname,
		user.Login,
		user.Password,
		user.Role,
		user.Position,
	).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("scanning user id: %w", err)
	}

	return userId, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users 
		SET name=$2, surname=$3, role=$4, position=$5
		WHERE id = $1`,
		user.ID,
		user.Name,
		user.Surname,
		user.Role,
		user.Position,
	)
	if err != nil {
		return fmt.Errorf("updating user %d: %w", user.ID, err)
	}

	return nil
}

func (r *UserRepository) MarkUserDisabled(ctx context.Context, id int16) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET disabled_at=now() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting user %d: %w", id, err)
	}

	return nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, `SELECT id, name, surname, login, password, created_at, role, position
                                 FROM users
                                 WHERE lower(login) = lower($1)`, login).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.Role,
		&user.Position,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting user by login: %w", err)
	}

	return &user, nil
}
