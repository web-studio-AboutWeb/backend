package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"web-studio-backend/internal/app/repository/database"
)

type client struct {
	conn *pgxpool.Pool
}

// NewClient connects to PostgreSQL and returns database.Database interface implementation.
func NewClient(ctx context.Context, connString string) (database.Database, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connectiong to postgresql: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging postgresql: %w", err)
	}

	return &client{conn: pool}, nil
}

func (c *client) Close() {
	c.conn.Close()
}
