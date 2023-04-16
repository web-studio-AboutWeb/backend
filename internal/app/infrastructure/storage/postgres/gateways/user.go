package gateways

import (
	"context"
	"errors"
	"fmt"
	user_core "web-studio-backend/internal/app/core/user"
	user_dto "web-studio-backend/internal/app/core/user/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
	"github.com/jackc/pgx/v5"
)

type UserGateway interface {
	CreateUser(ctx context.Context, user *user_core.User) (int16, error)
	GetUser(ctx context.Context, dto *user_dto.UserGet) (*user_core.User, error)
	GetUserByLogin(ctx context.Context, login string) (*user_core.User, error)
	UpdateUser(ctx context.Context, user *user_core.User) error
	DeleteUser(ctx context.Context, dto *user_dto.UserDelete) error
}

type userGateway struct {
	client *postgres.Client
}

func NewUserGateway(client *postgres.Client) UserGateway {
	return &userGateway{client: client}
}

func (c *userGateway) CreateUser(
	ctx context.Context, user *user_core.User,
) (int16, error) {
	row := c.client.Conn.QueryRow(ctx,
		`insert into users(name, surname, login, password, role)
             values($1, $2, $3, $4, $5)
             returning  id`,
		user.Name,
		user.Surname,
		user.Login,
		user.Password,
		user.Role,
	)

	var userId int16
	if err := row.Scan(&userId); err != nil {
		return 0, fmt.Errorf("scanning user id: %w", err)
	}

	return userId, nil
}

func (c *userGateway) GetUser(ctx context.Context, dto *user_dto.UserGet) (*user_core.User, error) {
	row := c.client.Conn.QueryRow(ctx, `select id, name, surname, created_at, role
                                 from users
                                 where id = $1`, dto.UserId)

	var user user_core.User
	if err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.CreatedAt,
		&user.Role,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, postgres.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting user %d: %w", dto.UserId, err)
	}

	return &user, nil
}

func (c *userGateway) GetUserByLogin(ctx context.Context, login string) (*user_core.User, error) {
	row := c.client.Conn.QueryRow(ctx, `select id, name, surname, login, password, created_at, role
                                 from users
                                 where lower(login) = lower($1)`, login)

	var user user_core.User
	if err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.Role,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, postgres.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting user by login %s: %w", login, err)
	}

	return &user, nil
}

func (c *userGateway) UpdateUser(ctx context.Context, user *user_core.User) error {
	_, err := c.client.Conn.Exec(ctx, `update users set name = $2, surname = $3, role = $4  where id = $1`,
		user.Id,
		user.Name,
		user.Surname,
		user.Role,
	)
	if err != nil {
		return fmt.Errorf("updating user %d: %w", user.Id, err)
	}

	return nil
}

func (c *userGateway) DeleteUser(ctx context.Context, dto *user_dto.UserDelete) error {
	_, err := c.client.Conn.Exec(ctx, `delete from users where id = $1`, dto.UserId)
	if err != nil {
		return fmt.Errorf("deleting user %d: %w", dto.UserId, err)
	}

	return nil
}
