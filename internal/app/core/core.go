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
	GetUser(ctx context.Context, req *domain.GetUserRequest) (*domain.GetUserResponse, error)
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.CreateUserResponse, error)
	UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *domain.DeleteUserRequest) (*domain.DeleteUserResponse, error)

	GetProject(ctx context.Context, req *domain.GetProjectRequest) (*domain.GetProjectResponse, error)
	CreateProject(ctx context.Context, req *domain.CreateProjectRequest) (*domain.CreateProjectResponse, error)
	UpdateProject(ctx context.Context, req *domain.UpdateProjectRequest) (*domain.UpdateProjectResponse, error)
	DeleteProject(ctx context.Context, req *domain.DeleteProjectRequest) (*domain.DeleteProjectResponse, error)
	GetProjectParticipants(ctx context.Context, req *domain.GetProjectParticipantsRequest) (*domain.GetProjectParticipantsResponse, error)
}

// core implements Core interface.
type core struct {
	repo database.Database
}

// New returns Core instance.
func New(ctx context.Context) (Core, error) {
	db, err := postgres.NewClient(ctx, os.Getenv("DATABASE_DSN"))
	if err != nil {
		return nil, fmt.Errorf("creating postgres driver: %w", err)
	}

	return &core{repo: db}, nil
}
