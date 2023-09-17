package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository"

	"github.com/jackc/pgx/v5"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool}
}

func (r *ProjectRepository) GetProject(ctx context.Context, id int16) (*domain.Project, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, title, description, cover_id, started_at, ended_at, link
                                 FROM projects
                                 WHERE id = $1`, id)

	var project domain.Project
	if err := row.Scan(
		&project.ID,
		&project.Title,
		&project.Description,
		&project.CoverId,
		&project.StartedAt,
		&project.EndedAt,
		&project.Link,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting project %d: %w", id, err)
	}

	return &project, nil
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project *domain.Project) (int16, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO projects(title, description, started_at, ended_at, link)
             VALUES($1, $2, $3, $4, $5)
             RETURNING  id`,
		project.Title,
		project.Description,
		project.StartedAt,
		project.EndedAt,
		project.Link,
	)

	var projectId int16
	if err := row.Scan(&projectId); err != nil {
		return 0, fmt.Errorf("scanning project id: %w", err)
	}

	return projectId, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, project *domain.Project) error {
	_, err := r.pool.Exec(ctx, `UPDATE projects SET title=$2, description=$3, started_at=$4, ended_at=$5, link=$6 WHERE id = $1`,
		project.ID,
		project.Title,
		project.Description,
		project.StartedAt,
		project.EndedAt,
		project.Link,
	)
	if err != nil {
		return fmt.Errorf("updating project %d: %w", project.ID, err)
	}

	return nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id int16) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting project %d: %w", id, err)
	}

	return nil
}

func (r *ProjectRepository) GetProjectParticipants(ctx context.Context, projectID int16) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.name, u.surname, u.created_at, u.role, u.position
	 	FROM projects p
			JOIN project_participants pp ON pp.project_id = p.id
			JOIN users u ON u.id = pp.user_id
	 	WHERE p.id = $1`, projectID)
	if err != nil {
		return nil, fmt.Errorf("selectiong project %d participants: %w", projectID, err)
	}
	defer rows.Close()

	var (
		participant  domain.User
		participants []domain.User
	)
	for rows.Next() {
		if err := rows.Scan(
			&participant.ID,
			&participant.Name,
			&participant.Surname,
			&participant.CreatedAt,
			&participant.Role,
			&participant.Position,
		); err != nil {
			return nil, fmt.Errorf("scanning participant: %w", err)
		}

		participants = append(participants, participant)
	}

	return participants, nil
}
