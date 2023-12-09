package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperr"
	"web-studio-backend/internal/app/infrastructure/repository"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

//go:generate mockgen -source=project.go -destination=./mocks/project.go -package=mocks
type ProjectRepository interface {
	GetProject(ctx context.Context, id int32) (*domain.Project, error)
	GetProjects(ctx context.Context) ([]domain.Project, error)
	CreateProject(ctx context.Context, project *domain.Project) (int32, error)
	UpdateProject(ctx context.Context, project *domain.Project) error
	DisableProject(ctx context.Context, id int32) error

	GetParticipants(ctx context.Context, projectID int32) ([]domain.ProjectParticipant, error)
	GetParticipant(ctx context.Context, participantID, projectID int32) (*domain.ProjectParticipant, error)
	AddParticipant(ctx context.Context, participant *domain.ProjectParticipant) error
	UpdateParticipant(ctx context.Context, participant *domain.ProjectParticipant) error
	RemoveParticipant(ctx context.Context, participantID, projectID int32) error
	SetProjectImageID(ctx context.Context, projectID int32, imageID string) error
}

type ProjectService struct {
	filesDir    string
	projectRepo ProjectRepository
	userRepo    UserRepository
	teamRepo    TeamRepository
	fileRepo    FileRepository
}

func NewProjectService(repo ProjectRepository, userRepo UserRepository, teamRepo TeamRepository, fileRepo FileRepository) *ProjectService {
	return &ProjectService{"projects", repo, userRepo, teamRepo, fileRepo}
}

func (s *ProjectService) GetProject(ctx context.Context, id int32) (*domain.Project, error) {
	project, err := s.projectRepo.GetProject(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project %d: %w", id, err)
	}

	return project, nil
}

func (s *ProjectService) GetProjects(ctx context.Context) ([]domain.Project, error) {
	projects, err := s.projectRepo.GetProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting projects: %w", err)
	}

	return projects, nil
}

func (s *ProjectService) CreateProject(ctx context.Context, project *domain.Project) (*domain.Project, error) {
	if err := project.Validate(); err != nil {
		return nil, fmt.Errorf("validating project: %w", err)
	}

	if project.TeamID != nil {
		_, err := s.teamRepo.GetTeam(ctx, *project.TeamID)
		if err != nil {
			if errors.Is(err, repository.ErrObjectNotFound) {
				return nil, apperr.NewNotFound("team_id")
			}
			return nil, fmt.Errorf("getting team %d: %w", *project.TeamID, err)
		}
	}

	projectId, err := s.projectRepo.CreateProject(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("creating project: %w", err)
	}

	createdProject, err := s.projectRepo.GetProject(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("getting project %d: %w", projectId, err)
	}

	return createdProject, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, project *domain.Project) (*domain.Project, error) {
	if err := project.Validate(); err != nil {
		return nil, fmt.Errorf("validating project: %w", err)
	}

	if project.TeamID != nil {
		_, err := s.teamRepo.GetTeam(ctx, *project.TeamID)
		if err != nil {
			if errors.Is(err, repository.ErrObjectNotFound) {
				return nil, apperr.NewNotFound("team_id")
			}
			return nil, fmt.Errorf("getting team %d: %w", *project.TeamID, err)
		}
	}

	_, err := s.projectRepo.GetProject(ctx, project.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project %d before update: %w", project.ID, err)
	}

	err = s.projectRepo.UpdateProject(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("updating project %d: %w", project.ID, err)
	}

	updatedProject, err := s.projectRepo.GetProject(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("getting project %d after update: %w", project.ID, err)
	}

	return updatedProject, nil
}

func (s *ProjectService) GetParticipants(ctx context.Context, projectID int32) ([]domain.ProjectParticipant, error) {
	project, err := s.projectRepo.GetProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project %d: %w", project.ID, err)
	}

	participants, err := s.projectRepo.GetParticipants(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("getting project %d participants: %w", projectID, err)
	}

	return participants, nil
}

func (s *ProjectService) GetParticipant(ctx context.Context, participantID, projectID int32) (*domain.ProjectParticipant, error) {
	participant, err := s.projectRepo.GetParticipant(ctx, participantID, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("user_id")
		}
		return nil, fmt.Errorf("getting participant: %w", err)
	}

	return participant, nil
}

func (s *ProjectService) AddParticipant(ctx context.Context, participant *domain.ProjectParticipant) (*domain.ProjectParticipant, error) {
	if err := participant.Validate(); err != nil {
		return nil, fmt.Errorf("validating participant: %w", err)
	}

	_, err := s.userRepo.GetActiveUser(ctx, participant.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("user_id")
		}
		return nil, fmt.Errorf("getting active user: %w", err)
	}

	_, err = s.projectRepo.GetProject(ctx, participant.ProjectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project: %w", err)
	}

	_, err = s.projectRepo.GetParticipant(ctx, participant.UserID, participant.ProjectID)
	if err != nil && !errors.Is(err, repository.ErrObjectNotFound) {
		return nil, fmt.Errorf("getting project participant: %w", err)
	}
	if err == nil {
		return nil, apperr.NewDuplicate("User already in participants list.", "user_id")
	}

	err = s.projectRepo.AddParticipant(ctx, participant)
	if err != nil {
		return nil, fmt.Errorf("adding participant: %w", err)
	}

	addedParticipant, err := s.projectRepo.GetParticipant(ctx, participant.UserID, participant.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("getting created paricipant: %w", err)
	}

	return addedParticipant, nil
}

func (s *ProjectService) UpdateParticipant(ctx context.Context, participant *domain.ProjectParticipant) (*domain.ProjectParticipant, error) {
	if err := participant.Validate(); err != nil {
		return nil, fmt.Errorf("validating participant: %w", err)
	}

	_, err := s.projectRepo.GetParticipant(ctx, participant.UserID, participant.ProjectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("user_id")
		}
		return nil, fmt.Errorf("getting participant: %w", err)
	}

	err = s.projectRepo.UpdateParticipant(ctx, participant)
	if err != nil {
		return nil, fmt.Errorf("updating participant: %w", err)
	}

	updatedParticipant, err := s.projectRepo.GetParticipant(ctx, participant.UserID, participant.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("getting updated participant: %w", err)
	}

	return updatedParticipant, nil
}

func (s *ProjectService) RemoveParticipant(ctx context.Context, participantID, projectID int32) error {
	_, err := s.projectRepo.GetParticipant(ctx, participantID, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperr.NewNotFound("user_id")
		}
		return fmt.Errorf("getting participant: %w", err)
	}

	err = s.projectRepo.RemoveParticipant(ctx, participantID, projectID)
	if err != nil {
		return fmt.Errorf("removing participant: %w", err)
	}

	return nil
}

func (s *ProjectService) SetProjectImage(ctx context.Context, projectID int32, img []byte) error {
	if len(img) > 5<<20 {
		return apperr.NewInvalidRequest("Image is too big.", "file")
	}

	mt := mimetype.Detect(img)
	if !mt.Is("image/jpeg") &&
		!mt.Is("image/png") &&
		!mt.Is("image/webp") {
		return apperr.NewInvalidRequest("Invalid image mime type.", "file")
	}

	fileID := uuid.New().String()
	fileName := fileID + mt.Extension()

	project, err := s.projectRepo.GetProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperr.NewNotFound("project_id")
		}
		return fmt.Errorf("getting project %d: %w", projectID, err)
	}

	err = s.projectRepo.SetProjectImageID(ctx, projectID, fileName)
	if err != nil {
		return fmt.Errorf("setting project %d image: %w", projectID, err)
	}

	err = s.fileRepo.Save(ctx, img, filepath.Join(s.filesDir, fileName))
	if err != nil {
		return fmt.Errorf("saving project image: %w", err)
	}

	if project.ImageId != "" {
		err = s.fileRepo.Delete(ctx, filepath.Join(s.filesDir, project.ImageId))
		if err != nil {
			slog.Error("Deleting old project image", slog.String("error", err.Error()))
		}
	}

	return nil
}

func (s *ProjectService) GetProjectImage(ctx context.Context, projectID int32) (*domain.Project, error) {
	project, err := s.projectRepo.GetProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("project_id")
		}
		return nil, fmt.Errorf("getting project %d: %w", projectID, err)
	}

	if project.ImageId == "" {
		return nil, apperr.NewNotFound("image_id")
	}

	data, err := s.fileRepo.Read(ctx, filepath.Join(s.filesDir, project.ImageId))
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("image_id")
		}
		return nil, fmt.Errorf("reading team image: %w", err)
	}

	project.ImageContent = data

	return project, nil
}
