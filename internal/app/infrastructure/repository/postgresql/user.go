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

func (r *UserRepository) GetUser(ctx context.Context, id int32) (*domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(ctx, `
		SELECT 
		    id, name, surname, username, email, created_at, updated_at, disabled_at, role, is_teamlead, image_id
        FROM users
        WHERE id = $1`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DisabledAt,
		&user.Role,
		&user.IsTeamLead,
		&user.ImageID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	user.RoleName = user.Role.String()

	return &user, nil
}

func (r *UserRepository) GetActiveUser(ctx context.Context, id int32) (*domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT
		    id, name, surname, username, email, created_at, updated_at, disabled_at, role, is_teamlead, image_id
        FROM users
        WHERE id = $1 AND disabled_at IS NULL`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DisabledAt,
		&user.Role,
		&user.IsTeamLead,
		&user.ImageID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}

	user.RoleName = user.Role.String()

	return &user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]domain.User, error) {

	rows, err := r.pool.Query(ctx, `
		SELECT 
		    id, name, surname, username, email, created_at, updated_at, disabled_at, role, is_teamlead, image_id
        FROM users
        WHERE disabled_at IS NULL`)
	if err != nil {
		return nil, fmt.Errorf("selecting users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User

		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DisabledAt,
			&user.Role,
			&user.IsTeamLead,
			&user.ImageID,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning user: %w", err)
		}

		user.RoleName = user.Role.String()

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	var userId int32

	err := r.pool.QueryRow(ctx,
		`INSERT INTO users(name, surname, username, email, encoded_password, salt, role, is_teamlead, image_id)
		 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING  id`,
		user.Name,
		user.Surname,
		user.Username,
		user.Email,
		user.EncodedPassword,
		user.Salt,
		user.Role,
		user.IsTeamLead,
		user.ImageID,
	).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("scanning user id: %w", err)
	}

	return userId, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users 
		SET name=$2, surname=$3, role=$4, updated_at=now()
		WHERE id = $1`,
		user.ID,
		user.Name,
		user.Surname,
		user.Role,
	)
	if err != nil {
		return fmt.Errorf("updating user %d: %w", user.ID, err)
	}

	return nil
}

func (r *UserRepository) DisableUser(ctx context.Context, id int32) error {
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
        WHERE (lower(username)=lower($1) OR lower(email)=lower($1))
          		AND disabled_at IS NULL`,
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

func (r *UserRepository) CheckUserUniqueness(ctx context.Context, username, email string) (*domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, username, email
        FROM users
        WHERE (lower(username)=lower($1) OR lower(email)=lower($2))
          		AND disabled_at IS NULL`,
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
