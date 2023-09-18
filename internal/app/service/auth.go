package service

import (
	"context"
	"errors"
	"fmt"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/app/infrastructure/repository"
	"web-studio-backend/internal/pkg/auth/session"
)

type AuthRepository interface {
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	CheckUserExists(ctx context.Context, id int16) error
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) SignIn(ctx context.Context, req *domain.SignInRequest) (*domain.SignInResponse, error) {
	user, err := s.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperror.NewInvalidRequest("Invalid credentials.", "")
		}
		return nil, fmt.Errorf("getting user by login: %w", err)
	}

	if !user.ComparePassword(req.Password) {
		return nil, apperror.NewInvalidRequest("Invalid credentials.", "")
	}

	sessionID, csrfToken, err := session.New(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generating session id: %w", err)
	}

	return &domain.SignInResponse{
		SessionID: sessionID,
		CSRFToken: csrfToken,
	}, nil
}

func (s *AuthService) SignOut(_ context.Context, sessionID string) {
	session.Delete(sessionID)
}

func (s *AuthService) CheckUserExists(ctx context.Context, id int16) error {
	return s.CheckUserExists(ctx, id)
}
