package handlers

import (
	"context"
	"errors"

	errcore "web-studio-backend/internal/app/core/shared/errors"
	staffer_dto "web-studio-backend/internal/app/core/staffer/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type GetStafferHandler struct {
	gateway gateways.StafferGateway
}

func NewGetStafferHandler(gateway gateways.StafferGateway) *GetStafferHandler {
	return &GetStafferHandler{gateway: gateway}
}

func (h *GetStafferHandler) Execute(
	ctx context.Context, dto *staffer_dto.StafferGet,
) (*staffer_dto.StafferObject, error) {
	staffer, err := h.gateway.GetStaffer(ctx, dto)
	if err != nil {
		if errors.Is(err, postgres.ErrObjectNotFound) {
			return nil, errcore.ProjectNotFoundError
		}
		return nil, err
	}

	return &staffer_dto.StafferObject{Staffer: staffer}, nil
}

type GetStaffersHandler struct {
	gateway gateways.StafferGateway
}

func NewGetStaffersHandler(gateway gateways.StafferGateway) *GetStaffersHandler {
	return &GetStaffersHandler{gateway: gateway}
}

func (h *GetStaffersHandler) Execute(
	ctx context.Context, dto *staffer_dto.StaffersGet,
) (*staffer_dto.StaffersObject, error) {
	staffers, err := h.gateway.GetStaffers(ctx, dto)
	if err != nil {
		if errors.Is(err, postgres.ErrObjectNotFound) {
			return nil, errcore.ProjectNotFoundError
		}
		return nil, err
	}

	return &staffer_dto.StaffersObject{Staffers: staffers}, nil
}
