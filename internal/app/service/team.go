package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/app/infrastructure/repository"
)

const teamsDir = "teams"

type TeamRepository interface {
	GetTeam(ctx context.Context, id int32) (*domain.Team, error)
	GetTeams(ctx context.Context) ([]domain.Team, error)
	CreateTeam(ctx context.Context, team *domain.Team) (int32, error)
	UpdateTeam(ctx context.Context, team *domain.Team) error
	SetTeamImageID(ctx context.Context, teamID int32, imageID string) error
	DisableTeam(ctx context.Context, teamID int32) error
	EnableTeam(ctx context.Context, teamID int32) error
	CheckTeamUniqueness(ctx context.Context, title string) (*domain.Team, error)
}

type TeamService struct {
	repo     TeamRepository
	fileRepo FileRepository
}

func NewTeamService(repo TeamRepository, fileRepo FileRepository) *TeamService {
	return &TeamService{repo, fileRepo}
}

func (s *TeamService) GetTeam(ctx context.Context, id int32) (*domain.Team, error) {
	team, err := s.repo.GetTeam(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("team_id")
		}
		return nil, fmt.Errorf("getting team %d: %w", id, err)
	}

	return team, nil
}

func (s *TeamService) GetTeams(ctx context.Context) ([]domain.Team, error) {
	teams, err := s.repo.GetTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting teams: %w", err)
	}

	return teams, nil
}

func (s *TeamService) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	if err := team.Validate(); err != nil {
		return nil, fmt.Errorf("validating team: %w", err)
	}

	foundTeam, err := s.repo.CheckTeamUniqueness(ctx, team.Title)
	if err != nil && !errors.Is(err, repository.ErrObjectNotFound) {
		return nil, fmt.Errorf("checking team uniqueness: %w", err)
	}
	if foundTeam != nil {
		return nil, apperror.NewDuplicate("Title already taken.", "title")
	}

	id, err := s.repo.CreateTeam(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("creating team: %w", err)
	}

	createdTeam, err := s.repo.GetTeam(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", id, err)
	}

	return createdTeam, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	if err := team.Validate(); err != nil {
		return nil, fmt.Errorf("validating team: %w", err)
	}

	_, err := s.repo.GetTeam(ctx, team.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("team_id")
		}
		return nil, fmt.Errorf("getting team %d: %w", team.ID, err)
	}

	if team.Title != "" {
		foundTeam, err := s.repo.CheckTeamUniqueness(ctx, team.Title)
		if err != nil && !errors.Is(err, repository.ErrObjectNotFound) {
			return nil, fmt.Errorf("checking team uniqueness: %w", err)
		}
		if foundTeam != nil {
			return nil, apperror.NewDuplicate("Title already taken.", "title")
		}
	}

	err = s.repo.UpdateTeam(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("creating team: %w", err)
	}

	updatedTeam, err := s.repo.GetTeam(ctx, team.ID)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", team.ID, err)
	}

	return updatedTeam, nil
}

func (s *TeamService) SetTeamImage(ctx context.Context, teamID int32, img []byte) error {
	if len(img) > 5<<20 {
		return apperror.NewInvalidRequest("Image is too big.", "file")
	}

	mt := mimetype.Detect(img)
	if !mt.Is("image/jpeg") &&
		!mt.Is("image/png") &&
		!mt.Is("image/webp") {
		return apperror.NewInvalidRequest("Invalid image mime type.", "file")
	}

	fileID := uuid.New().String()
	fileName := fileID + mt.Extension()

	team, err := s.repo.GetTeam(ctx, teamID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperror.NewNotFound("team_id")
		}
		return fmt.Errorf("getting team %d: %w", teamID, err)
	}

	err = s.repo.SetTeamImageID(ctx, teamID, fileName)
	if err != nil {
		return fmt.Errorf("setting team %d image: %w", teamID, err)
	}

	err = s.fileRepo.Save(ctx, img, filepath.Join(teamsDir, fileName))
	if err != nil {
		return fmt.Errorf("saving team image: %w", err)
	}

	if team.ImageID != "" {
		err = s.fileRepo.Delete(ctx, filepath.Join(teamsDir, team.ImageID))
		if err != nil {
			slog.Error("Deleting old team image", slog.String("error", err.Error()))
		}
	}

	return nil
}

func (s *TeamService) GetTeamImage(ctx context.Context, teamID int32) (*domain.Team, error) {
	team, err := s.repo.GetTeam(ctx, teamID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("team_id")
		}
		return nil, fmt.Errorf("getting team %d: %w", teamID, err)
	}

	if team.ImageID == "" {
		return nil, apperror.NewNotFound("image_id")
	}

	data, err := s.fileRepo.Read(ctx, filepath.Join(teamsDir, team.ImageID))
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("image_id")
		}
		return nil, fmt.Errorf("reading team image: %w", err)
	}

	team.ImageContent = data

	return team, nil
}

func (s *TeamService) DisableTeam(ctx context.Context, teamID int32) error {
	_, err := s.repo.GetTeam(ctx, teamID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperror.NewNotFound("team_id")
		}
		return fmt.Errorf("getting team %d: %w", teamID, err)
	}

	err = s.repo.DisableTeam(ctx, teamID)
	if err != nil {
		return fmt.Errorf("disabling team %d: %w", teamID, err)
	}

	return nil
}

func (s *TeamService) EnableTeam(ctx context.Context, teamID int32) error {
	_, err := s.repo.GetTeam(ctx, teamID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperror.NewNotFound("team_id")
		}
		return fmt.Errorf("getting team %d: %w", teamID, err)
	}

	err = s.repo.EnableTeam(ctx, teamID)
	if err != nil {
		return fmt.Errorf("enabling team %d: %w", teamID, err)
	}

	return nil
}
