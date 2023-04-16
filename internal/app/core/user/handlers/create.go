package handlers

import (
	"context"
	"errors"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	user_core "web-studio-backend/internal/app/core/user"
	user_dto "web-studio-backend/internal/app/core/user/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type CreateUserHandler struct {
	gateway gateways.UserGateway
}

func NewCreateUserHandler(gateway gateways.UserGateway) *CreateUserHandler {
	return &CreateUserHandler{gateway: gateway}
}

func (h *CreateUserHandler) Execute(
	ctx context.Context, dto *user_dto.UserCreate,
) (*user_dto.UserObject, error) {
	user, err := h.gateway.GetUserByLogin(ctx, dto.Login)
	if err != nil && !errors.Is(err, postgres.ErrObjectNotFound) {
		return nil, errcore.NewInternalError(err)
	}
	if user != nil {
		return nil, user_core.LoginAlreadyTakenError
	}

	user = &user_core.User{
		Name:      dto.Name,
		Surname:   dto.Surname,
		Login:     dto.Login,
		Password:  dto.Password,
		Role:      dto.Role,
	}

	userId, err := h.gateway.CreateUser(ctx, user)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	user, err = h.gateway.GetUser(ctx, &user_dto.UserGet{UserId: userId})
	if err != nil {
		return nil, err
	}

	return &user_dto.UserObject{User: user}, nil
}