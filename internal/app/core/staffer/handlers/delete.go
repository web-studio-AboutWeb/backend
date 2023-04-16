package handlers

import (
	"context"
	staffer_dto "web-studio-backend/internal/app/core/staffer/dto"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type DeleteStafferHandler struct {
	gateway gateways.StafferGateway
	getStafferHandler GetStafferHandler
}

func NewDeleteStafferHandler(
	gateway gateways.StafferGateway, getStafferHandler *GetStafferHandler,
) *DeleteStafferHandler {
	return &DeleteStafferHandler{
		gateway: gateway, getStafferHandler: *getStafferHandler,
	}
}

func (h *DeleteStafferHandler) Execute(
	ctx context.Context, dto *staffer_dto.StafferDelete,
) (interface{}, error) {
	_, err := h.getStafferHandler.Execute(
		ctx, &staffer_dto.StafferGet{StafferId: dto.StafferId},
	)
	if err != nil {
		return nil, err
	}

	err = h.gateway.DeleteStaffer(ctx, dto)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return nil, nil
}