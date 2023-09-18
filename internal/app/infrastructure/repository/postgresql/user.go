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
	var user domain.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, name, surname, username, created_at, updated_at, disabled_at, role, position
        FROM users
        WHERE id = $1`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DisabledAt,
		&user.Role,
		&user.Position,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	user.RoleName = user.Role.String()
	user.PositionName = user.Position.String()

	return &user, nil
}

func (r *UserRepository) GetActiveUser(ctx context.Context, id int16) (*domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, surname, username, email, created_at, updated_at, role, position
        FROM users
        WHERE id = $1 AND disabled_at IS NULL`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
		&user.Position,
	)
	if err != nil {
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
		`INSERT INTO users(name, surname, username, email, encoded_password, salt, role, position)
		 VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING  id`,
		user.Name,
		user.Surname,
		user.Username,
		user.Email,
		user.EncodedPassword,
		user.Salt,
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
		SET name=$2, surname=$3, role=$4, position=$5, updated_at=now()
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

func (r *UserRepository) DisableUser(ctx context.Context, id int16) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET disabled_at=now() WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("deleting user %d: %w", id, err)
	}

	return nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, username, email, encoded_password, salt
        FROM users
        WHERE (lower(username)=lower($1) OR lower(email)=lower($1)) AND disabled_at IS NULL`,
		login).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
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

func (r *UserRepository) CheckUsernameUniqueness(ctx context.Context, username, email string) (*domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, username, email
        FROM users
        WHERE (lower(username)=lower($1) OR lower(email)=lower($2)) AND disabled_at IS NULL`,
		username, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	return &user, nil
}
