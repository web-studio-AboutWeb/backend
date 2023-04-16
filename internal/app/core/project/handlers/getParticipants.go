package handlers

import (
	"context"
	project_dto "web-studio-backend/internal/app/core/project/dto"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type GetProjectParticipantsHandler struct {
	gateway gateways.ProjectGateway
	getProjectHandler *GetProjectHandler
}

func NewGetProjectParticipantsHandler(
	gateway gateways.ProjectGateway, getProjectHandler *GetProjectHandler,
) *GetProjectParticipantsHandler {
	return &GetProjectParticipantsHandler{
		gateway: gateway, getProjectHandler: getProjectHandler,
	}
}

func (h *GetProjectParticipantsHandler) Execute(
	ctx context.Context, dto *project_dto.ProjectStaffersGet,
) (*project_dto.ProjectStaffers, error) {
	_, err := h.getProjectHandler.Execute(
		ctx, &project_dto.ProjectGet{ProjectId: dto.ProjectId},
	)
	if err != nil {
		return nil, err
	}

	staffers, err := h.gateway.GetProjectStaffers(ctx, dto)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return &project_dto.ProjectStaffers{Staffers: staffers}, nil
}