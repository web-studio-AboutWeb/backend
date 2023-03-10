package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/errcore"
	"web-studio-backend/internal/app/infrastructure/config"
	"web-studio-backend/internal/app/repository/database"

	"golang.org/x/crypto/bcrypt"
)

func (c *core) SignIn(ctx context.Context, req *domain.SignInRequest) (*domain.SignInResponse, error) {
	user, err := c.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, database.ErrObjectNotFound) {
			return nil, errcore.InvalidCredentialsError
		}
		return nil, errcore.NewInternalError(err)
	}

	if !user.CheckPassword(req.Password) {
		return nil, errcore.InvalidCredentialsError
	}

	accessToken, err := c.CreateAccessToken(user)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}
	refreshToken, err := c.CreateRefreshToken(user)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return &domain.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *core) SignUp(ctx context.Context, req *domain.SignUpRequest) (*domain.SignUpResponse, error) {
	err := c.repo.CheckUsernameUniqueness(ctx, req.Username)
	if err != nil && !errors.Is(err, database.ErrObjectNotFound) {
		return nil, errcore.NewInternalError(err)
	}
	if err == nil {
		return nil, errcore.UsernameTakenError
	}

	err = c.repo.CheckEmailUniqueness(ctx, req.Email)
	if err != nil && !errors.Is(err, database.ErrObjectNotFound) {
		return nil, errcore.NewInternalError(err)
	}
	if err == nil {
		return nil, errcore.EmailTakenError
	}

	hashedPassword, err := c.HashPassword(req.Password)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}
	req.Password = hashedPassword

	id, err := c.repo.SignUp(ctx, req)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	user, err := c.repo.GetUser(ctx, &domain.GetUserRequest{UserId: id})
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	accessToken, err := c.CreateAccessToken(user)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}
	refreshToken, err := c.CreateRefreshToken(user)
	if err != nil {
		return nil, errcore.NewInternalError(err)
	}

	return &domain.SignUpResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *core) CreateAccessToken(user *domain.User) (string, error) {
	metadata := domain.AccessTokenMetadata{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	claims := jwt.MapClaims{
		"user": metadata,
		"exp":  time.Duration(config.Config.App.Jwt.AccessExpirationMinutes) * time.Minute,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := accessToken.SignedString([]byte(config.Config.App.Jwt.AccessTokenSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *core) CreateRefreshToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Duration(config.Config.App.Jwt.RefreshExpirationDays) * time.Hour * 24,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := refreshToken.SignedString([]byte(config.Config.App.Jwt.RefreshTokenSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *core) ValidateToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errcore.ErrInvalidToken
}

func (c *core) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
