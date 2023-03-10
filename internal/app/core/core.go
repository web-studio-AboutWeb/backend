package core

import (
	"context"
	"fmt"
	"os"
	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/repository/database"
	"web-studio-backend/internal/app/repository/database/postgres"
)

// Core represents business logic layer interface.
type Core interface {
	SignIn(ctx context.Context, req *domain.SignInRequest) (*domain.SignInResponse, error)
	SignUp(ctx context.Context, req *domain.SignUpRequest) (*domain.SignUpResponse, error)

	GetUser(ctx context.Context, req *domain.GetUserRequest) (*domain.GetUserResponse, error)
}

// core implements Core interface.
type core struct {
	repo database.Database
}

// New returns Core instance.
func New(ctx context.Context) (Core, error) {
	db, err := postgres.NewDriver(ctx, os.Getenv("DATABASE_DSN"))
	if err != nil {
		return nil, fmt.Errorf("creating postgres driver: %w", err)
	}

	return &core{repo: db}, nil
}
