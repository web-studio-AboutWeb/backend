package dto

import (
	"web-studio-backend/internal/app/domain"
)

type (
	CreateUserIn struct {
		Name       string          `json:"name"`
		Surname    string          `json:"surname"`
		Username   string          `json:"username"`
		Email      string          `json:"email"`
		Password   string          `json:"password"`
		Role       domain.UserRole `json:"role"`
		IsTeamLead bool            `json:"isTeamLead"`
	}

	UpdateUserIn struct {
		Name       string          `json:"name"`
		Surname    string          `json:"surname"`
		Role       domain.UserRole `json:"role"`
		IsTeamLead bool            `json:"isTeamLead"`
	}
)

func (in *CreateUserIn) ToDomain() *domain.User {
	if in == nil {
		return nil
	}

	return &domain.User{
		Name:            in.Name,
		Surname:         in.Surname,
		Username:        in.Username,
		Email:           in.Email,
		EncodedPassword: in.Password,
		Role:            in.Role,
		IsTeamLead:      in.IsTeamLead,
	}
}

func (in *UpdateUserIn) ToDomain(userID int32) *domain.User {
	if in == nil {
		return nil
	}

	return &domain.User{
		ID:         userID,
		Name:       in.Name,
		Surname:    in.Surname,
		Role:       in.Role,
		IsTeamLead: in.IsTeamLead,
	}
}
