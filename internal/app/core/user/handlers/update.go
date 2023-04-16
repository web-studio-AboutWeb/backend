package handlers

import (
	"context"
	"encoding/json"
	"web-studio-backend/internal/app/infrastructure/logger"
	errcore "web-studio-backend/internal/app/core/shared/errors"
	user_dto "web-studio-backend/internal/app/core/user/dto"
	"web-studio-backend/internal/app/infrastructure/storage/postgres/gateways"
)

type UpdateUserHandler struct {
	gateway gateways.UserGateway
	getUserHandler *GetUserHandler
}

func NewUpdateUserHandler(
	gateway gateways.UserGateway, getUserHandler *GetUserHandler,
) *UpdateUserHandler {
	return &UpdateUserHandler{gateway: gateway, getUserHandler: getUserHandler}
}

func (h *UpdateUserHandler) Execute(
	ctx context.Context, dto *user_dto.UserUpdate,
) (*user_dto.UserObject, error) {
	resp, err := h.getUserHandler.Execute(
		ctx, &user_dto.UserGet{UserId: dto.UserId},
	)
	if err != nil {
		return nil, err
	}

	user := resp.User
	user.Name = dto.Name
	user.Surname = dto.Surname
	user.Role = dto.Role

    jsonBytes, err := json.Marshal(user)
    if err != nil {
        panic(err)
    }
    userString := string(jsonBytes)
	logger.Logger.Info().Msg(userString)

	err = h.gateway.UpdateUser(ctx, user)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}


	response, err := h.getUserHandler.Execute(
		ctx, &user_dto.UserGet{UserId: dto.UserId},
	)
	if err != nil {
		return nil, err
	}

	return &user_dto.UserObject{User: response.User}, nil
}