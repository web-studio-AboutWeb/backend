package service

import (
	"context"
	"errors"
	"fmt"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperr"
	"web-studio-backend/internal/app/infrastructure/repository"
)

type ProjectCategoryRepository interface {
	CreateProjectCategory(ctx context.Context, pc *domain.ProjectCategory) (int16, error)
	GetProjectCategories(ctx context.Context) ([]domain.ProjectCategory, error)
	GetProjectCategory(ctx context.Context, id int16) (*domain.ProjectCategory, error)
	GetProjectCategoryByName(ctx context.Context, name string) (*domain.ProjectCategory, error)
	UpdateProjectCategory(ctx context.Context, pc *domain.ProjectCategory) error
	DeleteProjectCategory(ctx context.Context, id int16) error
}

type ProjectCategoryService struct {
	repo ProjectCategoryRepository
}

func NewProjectCategoryService(repo ProjectCategoryRepository) *ProjectCategoryService {
	return &ProjectCategoryService{repo}
}

func (s *ProjectCategoryService) CreateProjectCategory(ctx context.Context, pc *domain.ProjectCategory) (*domain.ProjectCategory, error) {
	// TODO: validate

	_, err := s.repo.GetProjectCategoryByName(ctx, pc.Name)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewDuplicate("Project category with such name already exists.", "name")
		}
		return nil, fmt.Errorf("getting project category by name %s: %w", pc.Name, err)
	}

	id, err := s.repo.CreateProjectCategory(ctx, pc)
	if err != nil {
		return nil, fmt.Errorf("creating project category: %w", err)
	}

	createdPc, err := s.repo.GetProjectCategory(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting created category: %w", err)
	}

	return createdPc, nil
}

func (s *ProjectCategoryService) GetProjectCategories(ctx context.Context) ([]domain.ProjectCategory, error) {
	return s.repo.GetProjectCategories(ctx)
}

func (s *ProjectCategoryService) UpdateProjectCategory(ctx context.Context, pc *domain.ProjectCategory) error {
	_, err := s.repo.GetProjectCategory(ctx, pc.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperr.NewNotFound("id")
		}
		return fmt.Errorf("getting project category %d: %w", pc.ID, err)
	}

	_, err = s.repo.GetProjectCategoryByName(ctx, pc.Name)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperr.NewDuplicate("Project category with such name already exists.", "name")
		}
		return fmt.Errorf("getting project category by name %s: %w", pc.Name, err)
	}

	err = s.repo.UpdateProjectCategory(ctx, pc)
	if err != nil {
		return fmt.Errorf("updating project category: %w", err)
	}

	return nil
}

func (s *ProjectCategoryService) DeleteProjectCategory(ctx context.Context, id int16) error {
	_, err := s.repo.GetProjectCategory(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperr.NewNotFound("id")
		}
		return fmt.Errorf("getting project category %d: %w", id, err)
	}

	err = s.repo.DeleteProjectCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting project category: %w", err)
	}

	return nil
}
