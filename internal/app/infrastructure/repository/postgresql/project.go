package postgresql

import (
	"context"
	"errors"
	"fmt"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository"

	"github.com/jackc/pgx/v5"
)

type ProjectRepository struct {
	pool Driver
}

func NewProjectRepository(dr Driver) *ProjectRepository {
	return &ProjectRepository{dr}
}

func (r *ProjectRepository) GetProject(ctx context.Context, id int32) (*domain.Project, error) {
	var project domain.Project

	err := r.pool.QueryRow(ctx, `
		SELECT id, title, description, cover_id, started_at, ended_at, link, isactive
        FROM projects
        WHERE id = $1`, id).Scan(
		&project.ID,
		&project.Title,
		&project.Description,
		&project.CoverId,
		&project.StartedAt,
		&project.EndedAt,
		&project.Link,
		&project.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning project: %w", err)
	}

	return &project, nil
}

func (r *ProjectRepository) GetActiveProject(ctx context.Context, id int32) (*domain.Project, error) {
	var project domain.Project

	err := r.pool.QueryRow(ctx, `
		SELECT id, title, description, cover_id, started_at, ended_at, link, isactive
        FROM projects
        WHERE id=$1 AND isactive`, id).Scan(
		&project.ID,
		&project.Title,
		&project.Description,
		&project.CoverId,
		&project.StartedAt,
		&project.EndedAt,
		&project.Link,
		&project.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning project: %w", err)
	}

	return &project, nil
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project *domain.Project) (int32, error) {
	var projectId int32

	err := r.pool.QueryRow(ctx,
		`INSERT INTO projects(title, description, started_at, link)
		 VALUES($1, $2, $3, $4)
		 RETURNING  id`,
		project.Title,
		project.Description,
		project.StartedAt,
		project.Link,
	).Scan(&projectId)
	if err != nil {
		return 0, fmt.Errorf("scanning project id: %w", err)
	}

	return projectId, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, project *domain.Project) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE projects
		SET title=$2, description=$3, started_at=$4, link=$5
		WHERE id = $1`,
		project.ID,
		project.Title,
		project.Description,
		project.StartedAt,
		project.Link,
	)
	if err != nil {
		return fmt.Errorf("updating project: %w", err)
	}

	return nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id int32) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}

	return nil
}

func (r *ProjectRepository) GetParticipants(ctx context.Context, projectID int32) ([]domain.ProjectParticipant, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
		    pp.user_id, pp.project_id, pp.role, pp.position, u.name, u.surname, u.username
	 	FROM project_participants pp
			JOIN users u ON u.id = pp.user_id
	 	WHERE pp.project_id = $1`, projectID)
	if err != nil {
		return nil, fmt.Errorf("selectiong project %d participants: %w", projectID, err)
	}
	defer rows.Close()

	var (
		p            domain.ProjectParticipant
		participants []domain.ProjectParticipant
	)
	for rows.Next() {
		if err = rows.Scan(
			&p.UserID,
			&p.ProjectID,
			&p.Role,
			&p.Position,
			&p.Name,
			&p.Surname,
			&p.Username,
		); err != nil {
			return nil, fmt.Errorf("scanning participant: %w", err)
		}

		participants = append(participants, p)
	}

	return participants, nil
}

func (r *ProjectRepository) GetParticipant(ctx context.Context, participantID, projectID int32) (*domain.ProjectParticipant, error) {
	var p *domain.ProjectParticipant

	err := r.pool.QueryRow(ctx, `
		SELECT 
		    pp.user_id, pp.project_id, pp.role, pp.position, u.name, u.surname, u.username
		FROM project_participants pp
			JOIN users u ON u.id=pp.user_id
		WHERE pp.user_id=$1 AND pp.project_id=$2`, participantID, projectID).Scan(
		&p.UserID,
		&p.ProjectID,
		&p.Role,
		&p.Position,
		&p.Name,
		&p.Surname,
		&p.Username,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning participant: %w", err)
	}

	return p, nil
}

func (r *ProjectRepository) AddParticipant(ctx context.Context, participant *domain.ProjectParticipant) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO project_participants(project_id, user_id, role, position) 
		VALUES ($1,$2,$3,$4)`,
		participant.ProjectID,
		participant.Username,
		participant.Role,
		participant.Position,
	)
	if err != nil {
		return fmt.Errorf("inserting participant: %w", err)
	}

	return nil
}

func (r *ProjectRepository) UpdateParticipant(ctx context.Context, participant *domain.ProjectParticipant) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE project_participants
		SET role=$3, position=$4
		WHERE user_id=$1 AND project_id=$2`,
		participant.UserID,
		participant.ProjectID,
		participant.Role,
		participant.Position,
	)
	if err != nil {
		return fmt.Errorf("updating participant: %w", err)
	}

	return nil
}

func (r *ProjectRepository) RemoveParticipant(ctx context.Context, participantID, projectID int32) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM project_participants
		WHERE user_id=$1 AND project_id=$2`, participantID, projectID)
	if err != nil {
		return fmt.Errorf("deleting participant: %w", err)
	}

	return nil
}
