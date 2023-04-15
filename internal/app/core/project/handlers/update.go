package handlers

import (
	"context"
	project_dto "web-studio-backend/internal/app/core/project/dto"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type UpdateProjectHandler struct {
	gateway gateways.ProjectGateway
	getProjectHandler GetProjectHandler
}

func NewUpdateProjectHandler(
	gateway gateways.ProjectGateway, getProjectHandler *GetProjectHandler,
) *UpdateProjectHandler {
	return &UpdateProjectHandler{
		gateway: gateway, getProjectHandler: *getProjectHandler,
	}
}

func (h *UpdateProjectHandler) Execute(
	ctx context.Context, dto *project_dto.ProjectUpdate,
) (*project_dto.ProjectObject, error) {
	resp, err := h.getProjectHandler.Execute(
		ctx, &project_dto.ProjectGet{ProjectId: dto.ProjectId},
	)
	if err != nil {
		return nil, err
	}
	project := resp.Project
	project.Id = dto.ProjectId
	project.Title = dto.Title
	project.Description = dto.Description
	project.StartedAt = dto.StartedAt
	project.EndedAt = dto.EndedAt
	project.Link = dto.Link

	err = h.gateway.UpdateProject(ctx, project)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	resp, err = h.getProjectHandler.Execute(
		ctx, &project_dto.ProjectGet{ProjectId: dto.ProjectId},
	)
	if err != nil {
		return nil, err
	}

	return &project_dto.ProjectObject{Project: resp.Project}, nil
}