package handlers

import (
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type StafferHandlers struct {
	CreateStafferHandler *CreateStafferHandler
	GetStafferHandler    *GetStafferHandler
	GetStaffersHandler   *GetStaffersHandler
	UpdateStafferHandler *UpdateStafferHandler
	DeleteStafferHandler *DeleteStafferHandler
}

func New(gateways *gateways.Gateways) (*StafferHandlers, error) {
	getStafferHandler := NewGetStafferHandler(gateways.StafferGateway)
	getStaffersHandler := NewGetStaffersHandler(gateways.StafferGateway)
	createStafferHandler := NewCreateStafferHandler(
		gateways.StafferGateway, getStafferHandler)
	updateStafferHandler := NewUpdateStafferHandler(
		gateways.StafferGateway, getStafferHandler,
	)
	deleteStafferHandler := NewDeleteStafferHandler(
		gateways.StafferGateway, getStafferHandler,
	)

	return &StafferHandlers{
		CreateStafferHandler: createStafferHandler,
		GetStafferHandler:    getStafferHandler,
		GetStaffersHandler:   getStaffersHandler,
		UpdateStafferHandler: updateStafferHandler,
		DeleteStafferHandler: deleteStafferHandler,
	}, nil
}
