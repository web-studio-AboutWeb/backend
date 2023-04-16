package handlers

import (
	"context"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	staffer_dto "web-studio-backend/internal/app/core/staffer/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type UpdateStafferHandler struct {
	gateway gateways.StafferGateway
	getStafferHandler GetStafferHandler
}

func NewUpdateStafferHandler(
	gateway gateways.StafferGateway, getProjectHandler *GetStafferHandler,
) *UpdateStafferHandler {
	return &UpdateStafferHandler{
		gateway: gateway, getStafferHandler: *getProjectHandler,
	}
}

func (h *UpdateStafferHandler) Execute(
	ctx context.Context, dto *staffer_dto.StafferUpdate,
) (*staffer_dto.StafferObject, error) {
	resp, err := h.getStafferHandler.Execute(
		ctx, &staffer_dto.StafferGet{StafferId: dto.StafferId},
	)
	if err != nil {
		return nil, err
	}
	staffer := resp.Staffer
	staffer.Id = dto.StafferId
	staffer.ProjectId = dto.ProjectId
	staffer.Position = dto.Position

	err = h.gateway.UpdateStaffer(ctx, staffer)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	resp, err = h.getStafferHandler.Execute(
		ctx, &staffer_dto.StafferGet{StafferId: dto.StafferId},
	)
	if err != nil {
		return nil, err
	}

	return &staffer_dto.StafferObject{Staffer: resp.Staffer}, nil
}