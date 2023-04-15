package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	Conn *pgxpool.Pool
}

// NewClient connects to PostgreSQL and returns database.Database interface implementation.
func NewClient(ctx context.Context, connString string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connectiong to postgresql: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging postgresql: %w", err)
	}
	return &Client{Conn: pool}, nil
}

func (c *Client) Close() {
	c.Conn.Close()
}
