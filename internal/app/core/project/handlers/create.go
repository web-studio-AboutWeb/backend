package handlers

import (
	"context"
	project_core "web-studio-backend/internal/app/core/project"
	project_dto "web-studio-backend/internal/app/core/project/dto"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type CreateProjectHandler struct {
	gateway gateways.ProjectGateway
	getProjectHandler *GetProjectHandler
}

func NewCreateProjectHandler(
	gateway gateways.ProjectGateway, getProjectHandler *GetProjectHandler,
) *CreateProjectHandler {
	return &CreateProjectHandler{
		gateway: gateway, getProjectHandler: getProjectHandler,
	}
}

func (h *CreateProjectHandler) Execute(
	ctx context.Context, dto *project_dto.ProjectCreate,
) (*project_dto.ProjectObject, error) {
	project := &project_core.Project{
		Title: dto.Title,
		Description: dto.Description,
		StartedAt: dto.StartedAt,
		EndedAt: dto.EndedAt,
		Link: dto.Link,
	}

	projectId, err := h.gateway.CreateProject(ctx, project)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	resp, err := h.getProjectHandler.Execute(
		ctx, &project_dto.ProjectGet{ProjectId: projectId},
	)
	if err != nil {
		return nil, err
	}

	return &project_dto.ProjectObject{Project: resp.Project}, nil
}