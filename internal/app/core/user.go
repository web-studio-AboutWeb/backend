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
