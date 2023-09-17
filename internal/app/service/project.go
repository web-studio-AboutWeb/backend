package service

import (
	"context"
	"errors"
	"fmt"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/app/infrastructure/repository"
)

type ProjectRepository interface {
	GetProject(ctx context.Context, id int16) (*domain.Project, error)
	CreateProject(ctx context.Context, project *domain.Project) (int16, error)
	UpdateProject(ctx context.Context, project *domain.Project) error
	GetProjectParticipants(ctx context.Context, id int16) ([]domain.User, error)
}

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{repo}
}

func (s *ProjectService) GetProject(ctx context.Context, id int16) (*domain.Project, error) {
	project, err := s.repo.GetProject(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project %d: %w", id, err)
	}

	return project, nil
}

func (s *ProjectService) CreateProject(ctx context.Context, project *domain.Project) (*domain.Project, error) {
	// TODO: validate fields

	projectId, err := s.repo.CreateProject(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("creating project: %w", err)
	}

	createdProject, err := s.repo.GetProject(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("getting project %d: %w", projectId, err)
	}

	return createdProject, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, project *domain.Project) (*domain.Project, error) {
	// TODO: validate fields

	project, err := s.repo.GetProject(ctx, project.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project %d before update: %w", project.ID, err)
	}

	err = s.repo.UpdateProject(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("updating project %d: %w", project.ID, err)
	}

	project, err = s.repo.GetProject(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("getting project %d after update: %w", project.ID, err)
	}

	return project, nil
}

func (s *ProjectService) GetProjectParticipants(ctx context.Context, projectID int16) ([]domain.User, error) {
	project, err := s.repo.GetProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project %d: %w", project.ID, err)
	}

	participants, err := s.repo.GetProjectParticipants(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("getting project %d participants: %w", projectID, err)
	}

	return participants, nil
}
