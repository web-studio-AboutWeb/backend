package gateways

import (
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
)

type Gateways struct {
	UserGateway UserGateway
	ProjectGateway ProjectGateway
	StafferGateway StafferGateway
}

func New(client *postgres.Client) (*Gateways, error) {
	userGateway := NewUserGateway(client)
	projectGateway := NewProjectGateway(client)
	stafferGateway := NewStafferGateway(client)

	return &Gateways{
		UserGateway: userGateway,
		ProjectGateway: projectGateway,
		StafferGateway: stafferGateway,
	}, nil
}