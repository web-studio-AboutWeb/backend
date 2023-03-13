package core

import (
	"context"
	"errors"
	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/errcore"
	"web-studio-backend/internal/app/repository/database"
)

func (c *core) GetUser(ctx context.Context, req *domain.GetUserRequest) (*domain.GetUserResponse, error) {
	user, err := c.repo.GetUser(ctx, req)
	if err != nil {
		if errors.Is(err, database.ErrObjectNotFound) {
			return nil, errcore.UserNotFoundError
		}
		return nil, err
	}

	return &domain.GetUserResponse{User: user}, nil
}

func (c *core) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	user, err := c.repo.GetUserByLogin(ctx, req.Login)
	if err != nil && !errors.Is(err, database.ErrObjectNotFound) {
		return nil, errcore.NewInternalError(err)
	}
	if user != nil {
		return nil, errcore.LoginAlreadyTakenError
	}

	userId, err := c.repo.CreateUser(ctx, req)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	response, err := c.GetUser(ctx, &domain.GetUserRequest{UserId: userId})
	if err != nil {
		return nil, err
	}

	return &domain.CreateUserResponse{User: response.User}, nil
}

func (c *core) UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	_, err := c.GetUser(ctx, &domain.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	err = c.repo.UpdateUser(ctx, req)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	response, err := c.GetUser(ctx, &domain.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	return &domain.UpdateUserResponse{User: response.User}, nil
}

func (c *core) DeleteUser(ctx context.Context, req *domain.DeleteUserRequest) (*domain.DeleteUserResponse, error) {
	_, err := c.GetUser(ctx, &domain.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	err = c.repo.DeleteUser(ctx, req)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return &domain.DeleteUserResponse{}, nil
}
