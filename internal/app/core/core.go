package core

import (
	"context"
	"fmt"
	"web-studio-backend/internal/app/infrastructure/config"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
	user_handlers "web-studio-backend/internal/app/core/user/handlers"
	project_handlers "web-studio-backend/internal/app/core/project/handlers"
	staffer_handlers "web-studio-backend/internal/app/core/staffer/handlers"
)

type Core struct {
	UserHandlers *user_handlers.UserHandlers
	ProjectHandlers *project_handlers.ProjectHandlers
	StafferHandlers *staffer_handlers.StafferHandlers
}

// New returns Core instance.
func New(ctx context.Context, config *config.Config) (*Core, error) {
	client, err := postgres.NewClient(ctx, config.DatabaseConnString)
	if err != nil {
		return nil, fmt.Errorf("creating postgres driver: %w", err)
	}
	
	gateways, err := gateways.New(client)
	if err != nil {
		return nil, fmt.Errorf("creating gateways: %w", err)
	}
	userHandlers, err := user_handlers.New(gateways)
	if err != nil {
		return nil, fmt.Errorf("creating user handlers: %w", err)
	}
	projectHandlers, err := project_handlers.New(gateways)
	if err != nil {
		return nil, fmt.Errorf("creating project handlers: %w", err)
	}
	staffer_handlers, err := staffer_handlers.New(gateways)

	return &Core{
		UserHandlers: userHandlers, 
		ProjectHandlers: projectHandlers,
		StafferHandlers: staffer_handlers,
	}, nil
}
