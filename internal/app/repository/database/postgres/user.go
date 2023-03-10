package postgres

import (
	"context"
	"errors"
	"fmt"
	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/repository/database"

	"github.com/jackc/pgx/v5"
)

func (d *driver) GetUser(ctx context.Context, req *domain.GetUserRequest) (*domain.User, error) {
	// TODO: fix
	row := d.conn.QueryRow(ctx, `select id, username, email, created_at, updated_at, exp, role, privilege 
                                 from users
                                 where id = $1`, req.UserId)

	var user domain.User
	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning user by credentials: %w", err)
	}

	return &user, nil
}

func (d *driver) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	// TODO: fix
	err := d.conn.QueryRow(ctx, `select id, username, email, password, created_at, updated_at, exp, role, privilege 
                                     from users
                                     where lower(username) = $1 or lower(email) = $1`, login).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning user by credentials: %w", err)
	}

	return &user, nil
}

func (d *driver) SignUp(ctx context.Context, req *domain.SignUpRequest) (int64, error) {
	var id int64
	// TODO: fix
	err := d.conn.QueryRow(ctx, "insert into users(username, email, password, role, privilege) values($1, $2, $3, $4, $5) returning id",
		req.Username,
		req.Email,
		req.Password,
		domain.UserRoleUser,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting user for sign up: %w", err)
	}

	return id, nil
}
