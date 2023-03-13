package postgres

import (
	"context"
	"errors"
	"fmt"
	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/repository/database"

	"github.com/jackc/pgx/v5"
)

func (c *client) GetProject(ctx context.Context, req *domain.GetProjectRequest) (*domain.Project, error) {
	row := c.conn.QueryRow(ctx, `select id, title, description, cover_id, started_at, ended_at, link
                                 from projects
                                 where id = $1`, req.ProjectId)

	var project domain.Project
	if err := row.Scan(
		&project.Id,
		&project.Title,
		&project.Description,
		&project.CoverId,
		&project.StartedAt,
		&project.EndedAt,
		&project.Link,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting project %d: %w", req.ProjectId, err)
	}

	return &project, nil
}

func (c *client) CreateProject(ctx context.Context, req *domain.CreateProjectRequest) (int16, error) {
	row := c.conn.QueryRow(ctx,
		`insert into projects(title, description, started_at, ended_at, link)
             values($1, $2, $3, $4, $5, $6)
             returning  id`,
		req.Title,
		req.Description,
		req.StartedAt,
		req.EndedAt,
		req.Link,
	)

	var projectId int16
	if err := row.Scan(&projectId); err != nil {
		return 0, fmt.Errorf("scanning project id: %w", err)
	}

	return projectId, nil
}

func (c *client) UpdateProject(ctx context.Context, req *domain.UpdateProjectRequest) error {
	_, err := c.conn.Exec(ctx, `update projects set title=$2, description=$3, started_at=$4, ended_at=$5, link=$6 where id = $1`,
		req.ProjectId,
		req.Title,
		req.Description,
		req.StartedAt,
		req.EndedAt,
		req.Link,
	)
	if err != nil {
		return fmt.Errorf("updating project %d: %w", req.ProjectId, err)
	}

	return nil
}

func (c *client) DeleteProject(ctx context.Context, req *domain.DeleteProjectRequest) error {
	_, err := c.conn.Exec(ctx, `delete from projects where id = $1`, req.ProjectId)
	if err != nil {
		return fmt.Errorf("deleting project %d: %w", req.ProjectId, err)
	}

	return nil
}

func (c *client) GetProjectParticipants(ctx context.Context, req *domain.GetProjectParticipantsRequest) ([]domain.User, error) {
	rows, err := c.conn.Query(ctx, `select u.id, u.name, u.surname, u.created_at, u.role, u.position
                                 from projects p
                                 	inner join project_participants pp on pp.project_id = p.id
                                 	inner join users u on u.id = pp.user_id
                                 where p.id = $1`, req.ProjectId)
	if err != nil {
		return nil, fmt.Errorf("selectiong project %d participants: %w", req.ProjectId, err)
	}
	defer rows.Close()

	var (
		participant  domain.User
		participants []domain.User
	)
	for rows.Next() {
		if err := rows.Scan(
			&participant.Id,
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
