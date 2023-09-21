package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperr"
	"web-studio-backend/internal/app/infrastructure/repository"
	"web-studio-backend/internal/pkg/strhelp"
)

//go:generate mockgen -source=user.go -destination=./mocks/user.go -package=mocks
type UserRepository interface {
	GetUser(ctx context.Context, id int32) (*domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	GetActiveUser(ctx context.Context, id int32) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (int32, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DisableUser(ctx context.Context, id int32) error
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	CheckUserUniqueness(ctx context.Context, username, email string) (*domain.User, error)
	SetUserImage(ctx context.Context, userID int32, imageID string) error
}

type UserService struct {
	filesDir string
	repo     UserRepository
	fileRepo FileRepository
}

func NewUserService(repo UserRepository, fileRepo FileRepository) *UserService {
	return &UserService{"users", repo, fileRepo}
}

func (s *UserService) GetUser(ctx context.Context, id int32) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("user_id")
		}
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}

	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting users: %w", err)
	}

	return users, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validating user: %w", err)
	}

	if user.Username == "" || len(user.Username) > 20 {
		return nil, apperr.NewInvalidRequest(
			fmt.Sprintf("Username cannot be empty and must not exceed %d characters.", 20),
			"username",
		)
	}

	if !strhelp.ValidateEmail(user.Email) {
		return nil, apperr.NewInvalidRequest("Email has invalid format.", "email")
	}

	if user.EncodedPassword == "" || len(user.EncodedPassword) > 20 {
		return nil, apperr.NewInvalidRequest(
			fmt.Sprintf("Password cannot be empty and must not exceed %d characters.", 20),
			"login",
		)
	}

	foundUser, err := s.repo.CheckUserUniqueness(ctx, user.Username, user.Email)
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

		return nil, apperr.NewDuplicate(
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
			return nil, apperr.NewNotFound("user_id")
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
			return apperr.NewNotFound("user_id")
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

func (s *UserService) SetUserImage(ctx context.Context, userID int32, img []byte) error {
	if len(img) > 5<<20 {
		return apperr.NewInvalidRequest("Image is too big.", "file")
	}

	mt := mimetype.Detect(img)
	if !mt.Is("image/jpeg") &&
		!mt.Is("image/png") &&
		!mt.Is("image/webp") {
		return apperr.NewInvalidRequest("Invalid image mime type.", "file")
	}

	fileID := uuid.New().String()
	fileName := fileID + mt.Extension()

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return apperr.NewNotFound("user_id")
		}
		return fmt.Errorf("getting user %d: %w", userID, err)
	}

	err = s.repo.SetUserImage(ctx, userID, fileName)
	if err != nil {
		return fmt.Errorf("setting user %d image: %w", userID, err)
	}

	err = s.fileRepo.Save(ctx, img, filepath.Join(s.filesDir, fileName))
	if err != nil {
		return fmt.Errorf("saving user image: %w", err)
	}

	if user.ImageID != "" {
		err = s.fileRepo.Delete(ctx, filepath.Join(s.filesDir, user.ImageID))
		if err != nil {
			slog.Error("Deleting old user image", slog.String("error", err.Error()))
		}
	}

	return nil
}

func (s *UserService) GetUserImage(ctx context.Context, userID int32) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("user_id")
		}
		return nil, fmt.Errorf("getting user %d: %w", userID, err)
	}

	if user.ImageID == "" {
		return nil, apperr.NewNotFound("image_id")
	}

	data, err := s.fileRepo.Read(ctx, filepath.Join(s.filesDir, user.ImageID))
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, apperr.NewNotFound("image_id")
		}
		return nil, fmt.Errorf("reading user image: %w", err)
	}

	user.ImageContent = data

	return user, nil
}
