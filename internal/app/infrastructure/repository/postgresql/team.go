package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository"
)

type TeamRepository struct {
	pool Driver
}

func NewTeamRepository(pool Driver) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) GetTeam(ctx context.Context, id int32) (*domain.Team, error) {
	var team domain.Team

	err := r.pool.QueryRow(ctx, `
		SELECT 
		    id, title, description, image_id, created_at, updated_at, disabled_at
		FROM teams
		WHERE id=$1`, id).Scan(
		&team.ID,
		&team.Title,
		&team.Description,
		&team.ImageID,
		&team.CreatedAt,
		&team.UpdatedAt,
		&team.DisabledAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning team: %w", err)
	}

	team.HasImage = team.ImageID != ""

	return &team, nil
}

func (r *TeamRepository) GetTeams(ctx context.Context) ([]domain.Team, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT 
		    id, title, description, image_id, created_at, updated_at, disabled_at
		FROM teams
		ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("selecting teams: %w", err)
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var team domain.Team

		err = rows.Scan(
			&team.ID,
			&team.Title,
			&team.Description,
			&team.ImageID,
			&team.CreatedAt,
			&team.UpdatedAt,
			&team.DisabledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning team: %w", err)
		}

		team.HasImage = team.ImageID != ""

		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *domain.Team) (int32, error) {
	var id int32

	err := r.pool.QueryRow(ctx, `
		INSERT INTO teams(title, description, image_id)
		VALUES ($1, $2, '')
		RETURNING id`, team.Title, team.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting team: %w", err)
	}

	return id, nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, team *domain.Team) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE teams
		SET title=$2, description=$3, updated_at=now()
		WHERE id=$1`,
		team.ID,
		team.Title,
		team.Description,
	)
	if err != nil {
		return fmt.Errorf("updating team: %w", err)
	}

	return nil
}

func (r *TeamRepository) SetTeamImageID(ctx context.Context, teamID int32, imageID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE teams
		SET image_id=$2, updated_at=now()
		WHERE id=$1`,
		teamID,
		imageID,
	)
	if err != nil {
		return fmt.Errorf("updating team: %w", err)
	}

	return nil
}

func (r *TeamRepository) DisableTeam(ctx context.Context, teamID int32) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE teams
		SET disabled_at=now(), updated_at=now()
		WHERE id=$1`,
		teamID,
	)
	if err != nil {
		return fmt.Errorf("updating team: %w", err)
	}

	return nil
}

func (r *TeamRepository) EnableTeam(ctx context.Context, teamID int32) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE teams
		SET disabled_at=NULL, updated_at=now()
		WHERE id=$1`,
		teamID,
	)
	if err != nil {
		return fmt.Errorf("updating team: %w", err)
	}

	return nil
}

func (r *TeamRepository) CheckTeamUniqueness(ctx context.Context, title string) (*domain.Team, error) {
	var team domain.Team

	err := r.pool.QueryRow(ctx, `
		SELECT 
		    id, title, description, image_id, created_at, updated_at, disabled_at
		FROM teams
		WHERE lower(title)=lower($1)`, title).Scan(
		&team.ID,
		&team.Title,
		&team.Description,
		&team.ImageID,
		&team.CreatedAt,
		&team.UpdatedAt,
		&team.DisabledAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning team: %w", err)
	}

	return &team, nil
}
