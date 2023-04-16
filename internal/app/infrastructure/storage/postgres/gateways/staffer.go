package gateways

import (
	"context"
	"errors"
	"fmt"
	staffer_core "web-studio-backend/internal/app/core/staffer"
	staffer_dto "web-studio-backend/internal/app/core/staffer/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
	"github.com/jackc/pgx/v5"
)

type StafferGateway interface {
	CreateStaffer(
		ctx context.Context, staffer *staffer_core.Staffer,
	) (int16, error)
	GetStaffer(
		ctx context.Context, dto *staffer_dto.StafferGet,
	) (*staffer_core.Staffer, error)
	UpdateStaffer(
		ctx context.Context, staffer *staffer_core.Staffer,
	) error
	DeleteStaffer(
		ctx context.Context, dto *staffer_dto.StafferDelete,
	) error
}

type stafferGateway struct {
	client *postgres.Client
}

func NewStafferGateway(client *postgres.Client) StafferGateway {
	return &stafferGateway{client: client}
}

func (g *stafferGateway) CreateStaffer(
	ctx context.Context, staffer *staffer_core.Staffer,
) (int16, error) {
	row := g.client.Conn.QueryRow(ctx,
		`insert into staffers(user_id, project_id, position)
             values($1, $2, $3)
             returning  id`,
		staffer.UserId,
		staffer.ProjectId,
		staffer.Position,
	)

	var stafferId int16
	if err := row.Scan(&stafferId); err != nil {
		return 0, fmt.Errorf("scanning user id: %w", err)
	}

	return stafferId, nil
}

func (g *stafferGateway) GetStaffer(
	ctx context.Context, dto *staffer_dto.StafferGet,
) (*staffer_core.Staffer, error) {
	row := g.client.Conn.QueryRow(ctx, 
		`select id, user_id, project_id, position
            from staffers
            where id = $1`, 
			dto.StafferId,
		)

	var staffer staffer_core.Staffer
	if err := row.Scan(
		&staffer.Id,
		&staffer.UserId,
		&staffer.ProjectId,
		&staffer.Position,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, postgres.ErrObjectNotFound
		}
		return nil, fmt.Errorf("getting user %d: %w", dto.StafferId, err)
	}

	return &staffer, nil
}

func (g *stafferGateway) UpdateStaffer(
	ctx context.Context, staffer *staffer_core.Staffer,
) error {
	_, err := g.client.Conn.Exec(ctx, 
		`update staffers set user_id = $2, project_id = $3, position = $4
			where id = $1`,
		staffer.Id,
		staffer.UserId,
		staffer.ProjectId,
		staffer.Position,
	)
	if err != nil {
		return fmt.Errorf("updating user %d: %w", staffer.Id, err)
	}

	return nil
}

func (g *stafferGateway) DeleteStaffer(
	ctx context.Context, dto *staffer_dto.StafferDelete,
) error {
	_, err := g.client.Conn.Exec(ctx, 
		`delete from staffers where id = $1`, 
		dto.StafferId,
	)
	if err != nil {
		return fmt.Errorf("deleting user %d: %w", dto.StafferId, err)
	}

	return nil
}
