package gateways

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	project_core "web-studio-backend/internal/app/core/project"
	project_dto "web-studio-backend/internal/app/core/project/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
)

type ProjectGateway interface {
	CreateProject(ctx context.Context, project *project_core.Project) (int16, error)
	GetProject(
		ctx context.Context, dto *project_dto.ProjectGet,
	) (*project_core.Project, error)
	GetProjects(ctx context.Context) ([]project_core.Project, error)
	GetProjectStaffers(
		ctx context.Context, dto *project_dto.ProjectStaffersGet,
	) ([]project_dto.ProjectStaffer, error)
	UpdateProject(ctx context.Context, project *project_core.Project) error
	DeleteProject(ctx context.Context, dto *project_dto.ProjectDelete) error
}

type projectGateway struct {
	client *postgres.Client
}

func NewProjectGateway(client *postgres.Client) ProjectGateway {
	return &projectGateway{client: client}
}

func (c *projectGateway) CreateProject(ctx context.Context, project *project_core.Project) (int16, error) {
	row := c.client.Conn.QueryRow(ctx,
		`insert into projects(title, description, started_at, ended_at, link)
             values($1, $2, $3, $4, $5)
             returning  id`,
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

func (c *projectGateway) GetProject(
	ctx context.Context, dto *project_dto.ProjectGet,
) (*project_core.Project, error) {
	row := c.client.Conn.QueryRow(ctx, `select id, title, description, cover_id, started_at, ended_at, link
                                 from projects
                                 where id = $1`, dto.ProjectId)

	var project project_core.Project

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
			return nil, postgres.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting project %d: %w", dto.ProjectId, err)
	}
	return &project, nil
}

func (c *projectGateway) GetProjects(ctx context.Context) ([]project_core.Project, error) {
	rows, err := c.client.Conn.Query(ctx, `select id, title, description, cover_id, started_at, ended_at, link
                                 from projects`)
	if err != nil {
		return nil, fmt.Errorf("getting projects: %w", err)
	}
	defer rows.Close()

	var (
		project  project_core.Project
		projects []project_core.Project
	)

	for rows.Next() {
		if err = rows.Scan(
			&project.Id,
			&project.Title,
			&project.Description,
			&project.CoverId,
			&project.StartedAt,
			&project.EndedAt,
			&project.Link,
		); err != nil {
			return nil, fmt.Errorf("scanning project: %w", err)
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (c *projectGateway) GetProjectStaffers(
	ctx context.Context, dto *project_dto.ProjectStaffersGet,
) ([]project_dto.ProjectStaffer, error) {
	rows, err := c.client.Conn.Query(ctx,
		`select s.id, s.project_id, s.position, u.id, u.name, 
				u.surname, u.created_at, u.role
			from projects p
				inner join staffers s on s.project_id = p.id
				inner join users u on u.id = s.user_id
			where p.id = $1`, dto.ProjectId)
	if err != nil {
		return nil, fmt.Errorf("selectiong project %d participants: %w", dto.ProjectId, err)
	}
	defer rows.Close()

	var (
		staffer  project_dto.ProjectStaffer
		staffers []project_dto.ProjectStaffer
	)
	for rows.Next() {
		if err := rows.Scan(
			&staffer.Id,
			&staffer.ProjectId,
			&staffer.Position,
			&staffer.User.Id,
			&staffer.User.Name,
			&staffer.User.Surname,
			&staffer.User.CreatedAt,
			&staffer.User.Role,
		); err != nil {
			return nil, fmt.Errorf("scanning participant: %w", err)
		}

		staffers = append(staffers, staffer)
	}

	return staffers, nil
}

func (c *projectGateway) UpdateProject(ctx context.Context, project *project_core.Project) error {
	_, err := c.client.Conn.Exec(ctx, `update projects set title=$2, description=$3, started_at=$4, ended_at=$5, link=$6 where id = $1`,
		project.Id,
		project.Title,
		project.Description,
		project.StartedAt,
		project.EndedAt,
		project.Link,
	)
	if err != nil {
		return fmt.Errorf("updating project %d: %w", project.Id, err)
	}

	return nil
}

func (c *projectGateway) DeleteProject(ctx context.Context, dto *project_dto.ProjectDelete) error {
	_, err := c.client.Conn.Exec(ctx, `delete from projects where id = $1`, dto.ProjectId)
	if err != nil {
		return fmt.Errorf("deleting project %d: %w", dto.ProjectId, err)
	}

	return nil
}
