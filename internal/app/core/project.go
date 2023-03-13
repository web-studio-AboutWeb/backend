package core

import (
	"context"
	"errors"
	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/errcore"
	"web-studio-backend/internal/app/repository/database"
)

func (c *core) GetProject(ctx context.Context, req *domain.GetProjectRequest) (*domain.GetProjectResponse, error) {
	Project, err := c.repo.GetProject(ctx, req)
	if err != nil {
		if errors.Is(err, database.ErrObjectNotFound) {
			return nil, errcore.ProjectNotFoundError
		}
		return nil, err
	}

	return &domain.GetProjectResponse{Project: Project}, nil
}

func (c *core) CreateProject(ctx context.Context, req *domain.CreateProjectRequest) (*domain.CreateProjectResponse, error) {
	projectId, err := c.repo.CreateProject(ctx, req)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	response, err := c.GetProject(ctx, &domain.GetProjectRequest{ProjectId: projectId})
	if err != nil {
		return nil, err
	}

	return &domain.CreateProjectResponse{Project: response.Project}, nil
}

func (c *core) UpdateProject(ctx context.Context, req *domain.UpdateProjectRequest) (*domain.UpdateProjectResponse, error) {
	_, err := c.GetProject(ctx, &domain.GetProjectRequest{ProjectId: req.ProjectId})
	if err != nil {
		return nil, err
	}

	err = c.repo.UpdateProject(ctx, req)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	response, err := c.GetProject(ctx, &domain.GetProjectRequest{ProjectId: req.ProjectId})
	if err != nil {
		return nil, err
	}

	return &domain.UpdateProjectResponse{Project: response.Project}, nil
}

func (c *core) DeleteProject(ctx context.Context, req *domain.DeleteProjectRequest) (*domain.DeleteProjectResponse, error) {
	_, err := c.GetProject(ctx, &domain.GetProjectRequest{ProjectId: req.ProjectId})
	if err != nil {
		return nil, err
	}

	err = c.repo.DeleteProject(ctx, req)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return &domain.DeleteProjectResponse{}, nil
}

func (c *core) GetProjectParticipants(ctx context.Context, req *domain.GetProjectParticipantsRequest) (*domain.GetProjectParticipantsResponse, error) {
	_, err := c.GetProject(ctx, &domain.GetProjectRequest{ProjectId: req.ProjectId})
	if err != nil {
		return nil, err
	}

	participants, err := c.repo.GetProjectParticipants(ctx, req)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return &domain.GetProjectParticipantsResponse{Participants: participants}, nil
}
