package database

import (
	"context"
	"web-studio-backend/internal/app/domain"
)

// Database represents database interface.
type Database interface {
	GetUser(ctx context.Context, req *domain.GetUserRequest) (*domain.User, error)
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (int16, error)
	UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) error
	DeleteUser(ctx context.Context, req *domain.DeleteUserRequest) error
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)

	GetProject(ctx context.Context, req *domain.GetProjectRequest) (*domain.Project, error)
	CreateProject(ctx context.Context, req *domain.CreateProjectRequest) (int16, error)
	UpdateProject(ctx context.Context, req *domain.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, req *domain.DeleteProjectRequest) error
	GetProjectParticipants(ctx context.Context, req *domain.GetProjectParticipantsRequest) ([]domain.User, error)

	Close()
}
