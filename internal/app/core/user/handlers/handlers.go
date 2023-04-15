package handlers

import (
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type UserHandlers struct {
	CreateUserHandler *CreateUserHandler
	GetUserHandler *GetUserHandler
	UpdateUserHandler *UpdateUserHandler
	DeleteUserHandler *DeleteUserHandler
}

func New(gateways *gateways.Gateways) (*UserHandlers, error) {
	createUserHandler := NewCreateUserHandler(gateways.UserGateway)
	getUserHandler := NewGetUserHandler(gateways.UserGateway)
	updateUserHandler := NewUpdateUserHandler(
		gateways.UserGateway, getUserHandler,
	)
	deleteUserHandler := NewDeleteUserHandler(
		gateways.UserGateway, getUserHandler,
	)

	return &UserHandlers{
		CreateUserHandler: createUserHandler,
		GetUserHandler: getUserHandler,
		UpdateUserHandler: updateUserHandler,
		DeleteUserHandler: deleteUserHandler,
	}, nil
}