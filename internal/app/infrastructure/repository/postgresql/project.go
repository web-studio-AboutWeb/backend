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

func NewProjectRepository(pool Driver) *ProjectRepository {
	return &ProjectRepository{pool}
}

func (r *ProjectRepository) GetProject(ctx context.Context, id int32) (*domain.Project, error) {
	var project domain.Project

	err := r.pool.QueryRow(ctx, `
		SELECT 
		   p.id, title, description, image_id, created_at, updated_at, started_at, ended_at,
		   link, isactive, technologies, team_id, pc.name
       	FROM projects p
		JOIN project_categories pc ON p.category_id = pc.id
       	WHERE p.id = $1`, id).Scan(
		&project.ID,
		&project.Title,
		&project.Description,
		&project.ImageId,
		&project.CreatedAt,
		&project.UpdatedAt,
		&project.StartedAt,
		&project.EndedAt,
		&project.Link,
		&project.IsActive,
		&project.Technologies,
		&project.TeamID,
		&project.Category,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("scanning project: %w", err)
	}

	return &project, nil
}

func (r *ProjectRepository) GetProjects(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT 
		    p.id, title, description, image_id, created_at, updated_at, started_at, ended_at,
		    link, isactive, technologies, team_id, pc.name
        FROM projects p
		JOIN project_categories pc ON p.category_id = pc.id
        WHERE isactive
        ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("selecting projects: %w", err)
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var project domain.Project

		err = rows.Scan(
			&project.ID,
			&project.Title,
			&project.Description,
			&project.ImageId,
			&project.CreatedAt,
			&project.UpdatedAt,
			&project.StartedAt,
			&project.EndedAt,
			&project.Link,
			&project.IsActive,
			&project.Technologies,
			&project.TeamID,
			&project.Category,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning project: %w", err)
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project *domain.Project) (int32, error) {
	var projectId int32

	err := r.pool.QueryRow(ctx, `
		INSERT INTO projects(title, description, team_id, isactive, link, technologies, image_id, started_at, ended_at, category_id)
		VALUES($1, $2, $3, TRUE, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		project.Title,
		project.Description,
		project.TeamID,
		project.Link,
		project.Technologies,
		project.ImageId,
		project.StartedAt,
		project.EndedAt,
		project.CategoryID,
	).Scan(&projectId)
	if err != nil {
		return 0, fmt.Errorf("scanning project id: %w", err)
	}

	return projectId, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, project *domain.Project) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE projects
		SET title=$2, description=$3, link=$4, technologies=$5, started_at=$6, ended_at=$7, updated_at=now(), category_id=$8
		WHERE id = $1`,
		project.ID,
		project.Title,
		project.Description,
		project.Link,
		project.Technologies,
		project.StartedAt,
		project.EndedAt,
		project.CategoryID,
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

func (r *ProjectRepository) DisableProject(ctx context.Context, id int32) error {
	_, err := r.pool.Exec(ctx, `UPDATE projects SET isactive=FALSE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("disabling project: %w", err)
	}

	return nil
}

func (r *ProjectRepository) GetParticipants(ctx context.Context, projectID int32) ([]domain.ProjectParticipant, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
		    pp.user_id, pp.project_id, pp.role, pp.position, pp.created_at, pp.updated_at,
		    u.name, u.surname, u.username
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
			&p.CreatedAt,
			&p.UpdatedAt,
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
	var p domain.ProjectParticipant

	err := r.pool.QueryRow(ctx, `
		SELECT 
		    pp.user_id, pp.project_id, pp.role, pp.position, pp.created_at, pp.updated_at,
		    u.name, u.surname, u.username
		FROM project_participants pp
			JOIN users u ON u.id=pp.user_id
		WHERE pp.user_id=$1 AND pp.project_id=$2`, participantID, projectID).Scan(
		&p.UserID,
		&p.ProjectID,
		&p.Role,
		&p.Position,
		&p.CreatedAt,
		&p.UpdatedAt,
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

	return &p, nil
}

func (r *ProjectRepository) AddParticipant(ctx context.Context, participant *domain.ProjectParticipant) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO project_participants(project_id, user_id, role, position) 
		VALUES ($1,$2,$3,$4)`,
		participant.ProjectID,
		participant.UserID,
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
		SET role=$3, position=$4, updated_at=now()
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

func (r *ProjectRepository) SetProjectImageID(ctx context.Context, projectID int32, imageID string) error {
	if _, err := r.pool.Exec(ctx, `
	UPDATE projects
	SET image_id=$2, updated_at=now()
	WHERE id=$1`,
		projectID,
		imageID,
	); err != nil {
		return fmt.Errorf("updating project: %w", err)
	}

	return nil
}
