package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"web-studio-backend/internal/app/repository/database"
)

func (d *driver) CheckUsernameUniqueness(ctx context.Context, username string) error {
	var temp int
	// TODO: fix
	err := d.conn.QueryRow(ctx, "select id from users where lower(username) = lower($1)", username).Scan(&temp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.ErrObjectNotFound
		}
		return err
	}

	return nil
}

func (d *driver) CheckEmailUniqueness(ctx context.Context, email string) error {
	var temp int
	// TODO: fix
	err := d.conn.QueryRow(ctx, "select id from users where lower(email) = lower($1)", email).Scan(&temp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.ErrObjectNotFound
		}
		return err
	}

	return nil
}
