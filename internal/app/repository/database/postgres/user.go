package postgres

import (
	"context"
	"errors"
	"fmt"
	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/repository/database"

	"github.com/jackc/pgx/v5"
)

func (c *client) GetUser(ctx context.Context, req *domain.GetUserRequest) (*domain.User, error) {
	row := c.conn.QueryRow(ctx, `select id, name, surname, created_at, role, position
                                 from users
                                 where id = $1`, req.UserId)

	var user domain.User
	if err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.CreatedAt,
		&user.Role,
		&user.Position,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting user %d: %w", req.UserId, err)
	}

	return &user, nil
}

func (c *client) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (int16, error) {
	row := c.conn.QueryRow(ctx,
		`insert into users(name, surname, login, password, role, position)
             values($1, $2, $3, $4, $5, $6)
             returning  id`,
		req.Name,
		req.Surname,
		req.Login,
		req.Password,
		req.Role,
		req.Position,
	)

	var userId int16
	if err := row.Scan(&userId); err != nil {
		return 0, fmt.Errorf("scanning user id: %w", err)
	}

	return userId, nil
}

func (c *client) UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) error {
	_, err := c.conn.Exec(ctx, `update users set name = $2, surname = $3, role = $4, position = $5 where id = $1`,
		req.UserId,
		req.Name,
		req.Surname,
		req.Role,
		req.Position,
	)
	if err != nil {
		return fmt.Errorf("updating user %d: %w", req.UserId, err)
	}

	return nil
}

func (c *client) DeleteUser(ctx context.Context, req *domain.DeleteUserRequest) error {
	_, err := c.conn.Exec(ctx, `delete from users where id = $1`, req.UserId)
	if err != nil {
		return fmt.Errorf("deleting user %d: %w", req.UserId, err)
	}

	return nil
}

func (c *client) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	row := c.conn.QueryRow(ctx, `select id, name, surname, login, password, created_at, role, position
                                 from users
                                 where lower(login) = lower($1)`, login)

	var user domain.User
	if err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.Role,
		&user.Position,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting user by login %s: %w", login, err)
	}

	return &user, nil
}
