package database

import (
	"context"
	"web-studio-backend/internal/app/domain"
)

// Database represents database interface.
type Database interface {
	GetUser(ctx context.Context, req *domain.GetUserRequest) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	CheckUsernameUniqueness(ctx context.Context, username string) error
	CheckEmailUniqueness(ctx context.Context, email string) error

	SignUp(ctx context.Context, req *domain.SignUpRequest) (int64, error)

	Close()
}
