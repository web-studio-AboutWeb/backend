package handlers

import (
	"context"
	"errors"

	project_dto "web-studio-backend/internal/app/core/project/dto"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type GetProjectHandler struct {
	gateway gateways.ProjectGateway
}

func NewGetProjectHandler(gateway gateways.ProjectGateway) *GetProjectHandler {
	return &GetProjectHandler{gateway: gateway}
}

func (h *GetProjectHandler) Execute(
	ctx context.Context, dto *project_dto.ProjectGet,
) (*project_dto.ProjectObject, error) {
	project, err := h.gateway.GetProject(ctx, dto)
	if err != nil {
		if errors.Is(err, postgres.ErrObjectNotFound) {
			return nil, errcore.ProjectNotFoundError
		}
		return nil, err
	}

	return &project_dto.ProjectObject{Project: project}, nil
}

type GetProjectsHandler struct {
	gateway gateways.ProjectGateway
}

func NewGetProjectsHandler(gateway gateways.ProjectGateway) *GetProjectsHandler {
	return &GetProjectsHandler{gateway: gateway}
}

func (h *GetProjectsHandler) Execute(ctx context.Context) (*project_dto.ProjectsObject, error) {
	projects, err := h.gateway.GetProjects(ctx)
	if err != nil {
		if errors.Is(err, postgres.ErrObjectNotFound) {
			return nil, errcore.ProjectNotFoundError
		}
		return nil, err
	}

	return &project_dto.ProjectsObject{Projects: projects}, nil
}
