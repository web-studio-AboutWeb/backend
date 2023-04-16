package handlers

import (
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type StafferHandlers struct {
	CreateStafferHandler *CreateStafferHandler
	GetStafferHandler *GetStafferHandler
	UpdateStafferHandler *UpdateStafferHandler
	DeleteStafferHandler *DeleteStafferHandler
}

func New(gateways *gateways.Gateways) (*StafferHandlers, error) {
	getStafferHandler := NewGetStafferHandler(gateways.StafferGateway)
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
		GetStafferHandler: getStafferHandler,
		UpdateStafferHandler: updateStafferHandler,
		DeleteStafferHandler: deleteStafferHandler,
	}, nil
} 