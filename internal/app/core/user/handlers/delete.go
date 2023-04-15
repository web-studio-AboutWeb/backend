package handlers

import (
	"context"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	user_dto "web-studio-backend/internal/app/core/user/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type DeleteUserHandler struct {
	gateway gateways.UserGateway
	getUserHandler *GetUserHandler
}

func NewDeleteUserHandler(gateway gateways.UserGateway, getUserHandler *GetUserHandler) *DeleteUserHandler {
	return &DeleteUserHandler{gateway: gateway, getUserHandler: getUserHandler}
}

func (h *DeleteUserHandler) Execute(
	ctx context.Context, dto *user_dto.UserDelete,
) (interface{}, error) {
	_, err := h.getUserHandler.Execute(
		ctx, &user_dto.UserGet{UserId: dto.UserId},
	)
	if err != nil {
		return nil, err
	}

	err = h.gateway.DeleteUser(ctx, dto)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return nil, nil
}