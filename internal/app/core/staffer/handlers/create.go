package handlers

import (
	"context"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	staffer_core "web-studio-backend/internal/app/core/staffer"
	staffer_dto "web-studio-backend/internal/app/core/staffer/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type CreateStafferHandler struct {
	gateway gateways.StafferGateway
	getStafferHandler *GetStafferHandler
}

func NewCreateStafferHandler(
	gateway gateways.StafferGateway, getStafferHandler *GetStafferHandler,
) *CreateStafferHandler {
	return &CreateStafferHandler{
		gateway: gateway, getStafferHandler: getStafferHandler,
	}
}

func (h *CreateStafferHandler) Execute(
	ctx context.Context, dto *staffer_dto.StafferCreate,
) (*staffer_dto.StafferObject, error) {
	staffer := &staffer_core.Staffer{
		UserId: dto.UserId,
		ProjectId: dto.ProjectId,
		Position: dto.Position,
	}

	stafferId, err := h.gateway.CreateStaffer(ctx, staffer)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	resp, err := h.getStafferHandler.Execute(
		ctx, &staffer_dto.StafferGet{StafferId: stafferId},
	)
	if err != nil {
		return nil, err
	}

	return &staffer_dto.StafferObject{Staffer: resp.Staffer}, nil
}