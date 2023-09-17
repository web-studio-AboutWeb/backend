package service

import (
	"context"
	"errors"
	"fmt"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/app/infrastructure/repository"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int16) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (int16, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	MarkUserDisabled(ctx context.Context, id int16) error
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return nil
}

func (s *UserService) GetUser(ctx context.Context, id int16) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.UserNotFoundError
		}
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser, err := s.repo.GetUserByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, repository.ErrObjectNotFound) {
		return nil, fmt.Errorf("getting user by login %s: %w", err)
	}
	if foundUser != nil {
		return nil, apperror.LoginAlreadyTakenError
	}

	userId, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	createdUser, err := s.repo.GetUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("getting user %d: %w", userId, err)
	}

	return createdUser, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	_, err := s.repo.GetUser(ctx, user.Id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.UserNotFoundError
		}
		return nil, fmt.Errorf("getting user %d: %w", user.Id, err)
	}

	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("updating user %d: %w", user.Id, err)
	}

	updatedUser, err := s.repo.GetUser(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("updating user %d: %w", user.Id, err)
	}

	return updatedUser, nil
}

func (s *UserService) RemoveUser(ctx context.Context, id int16) error {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperror.UserNotFoundError
		}
		return fmt.Errorf("getting user %d: %w", id, err)
	}
	_ = user

	// TODO: compare user role

	err = s.repo.MarkUserDisabled(ctx, id)
	if err != nil {
		return fmt.Errorf("marking user %d disabled: %w", id, err)
	}

	return nil
}