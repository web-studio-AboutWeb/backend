package handlers

import (
	"context"
	project_dto "web-studio-backend/internal/app/core/project/dto"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type DeleteProjectHandler struct {
	gateway gateways.ProjectGateway
	getProjectHandler GetProjectHandler
}

func NewDeleteProjectHandler(
	gateway gateways.ProjectGateway, getProjectHandler *GetProjectHandler,
) *DeleteProjectHandler {
	return &DeleteProjectHandler{
		gateway: gateway, getProjectHandler: *getProjectHandler,
	}
}

func (h *DeleteProjectHandler) Execute(
	ctx context.Context, dto *project_dto.ProjectDelete,
) (interface{}, error) {
	_, err := h.getProjectHandler.Execute(
		ctx, &project_dto.ProjectGet{ProjectId: dto.ProjectId},
	)
	if err != nil {
		return nil, err
	}

	err = h.gateway.DeleteProject(ctx, dto)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return nil, nil
}