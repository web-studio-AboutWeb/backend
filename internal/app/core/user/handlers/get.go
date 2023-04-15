package handlers

import (
	"context"
	"errors"
	user_core "web-studio-backend/internal/app/core/user"
	user_dto "web-studio-backend/internal/app/core/user/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type GetUserHandler struct {
	gateway gateways.UserGateway
}

func NewGetUserHandler(gateway gateways.UserGateway) *GetUserHandler {
	return &GetUserHandler{gateway: gateway}
}

func (h *GetUserHandler) Execute(
	ctx context.Context, dto *user_dto.UserGet,
) (*user_dto.UserObject, error) {
	user, err := h.gateway.GetUser(ctx, dto)
	if err != nil {
		if errors.Is(err, postgres.ErrObjectNotFound) {
			return nil, user_core.UserNotFoundError
		}
		return nil, err
	}

	return &user_dto.UserObject{User: user}, nil
}