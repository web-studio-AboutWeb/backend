package service

import (
	"context"
	"errors"
	"fmt"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/app/infrastructure/repository"
)

//go:generate mockgen -source=user.go -destination=./mocks/user.go -package=mocks
type UserRepository interface {
	GetUser(ctx context.Context, id int32) (*domain.User, error)
	GetActiveUser(ctx context.Context, id int32) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (int32, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DisableUser(ctx context.Context, id int32) error
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	CheckUsernameUniqueness(ctx context.Context, username, email string) (*domain.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) GetUser(ctx context.Context, id int32) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("user_id")
		}
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validating user: %w", err)
	}

	foundUser, err := s.repo.CheckUsernameUniqueness(ctx, user.Username, user.Email)
	if err != nil && !errors.Is(err, repository.ErrObjectNotFound) {
		return nil, fmt.Errorf("getting user by login: %w", err)
	}
	if foundUser != nil {
		var field, msgField string
		if foundUser.Email == user.Email {
			field = "email"
			msgField = "Email"
		} else {
			field = "username"
			msgField = "Username"
		}

		return nil, apperror.NewDuplicate(
			fmt.Sprintf("%s already taken.", msgField),
			field,
		)
	}

	err = user.EncodePassword()
	if err != nil {
		return nil, fmt.Errorf("encoding user password: %w", err)
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
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validating user: %w", err)
	}

	_, err := s.repo.GetUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewNotFound("user_id")
		}
		return nil, fmt.Errorf("getting user %d: %w", user.ID, err)
	}

	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("updating user %d: %w", user.ID, err)
	}

	updatedUser, err := s.repo.GetUser(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("updating user %d: %w", user.ID, err)
	}

	return updatedUser, nil
}

func (s *UserService) RemoveUser(ctx context.Context, id int32) error {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperror.NewNotFound("user_id")
		}
		return fmt.Errorf("getting user %d: %w", id, err)
	}
	_ = user

	// TODO: compare user role

	err = s.repo.DisableUser(ctx, id)
	if err != nil {
		return fmt.Errorf("marking user %d disabled: %w", id, err)
	}

	return nil
}
