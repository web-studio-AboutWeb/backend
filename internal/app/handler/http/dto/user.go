package dto

import "web-studio-backend/internal/app/domain"

type (
	CreateUserIn struct {
		Name     string              `json:"name"`
		Surname  string              `json:"surname"`
		Login    string              `json:"login"`
		Password string              `json:"password"`
		Role     domain.UserRole     `json:"role"`
		Position domain.UserPosition `json:"position"`
	}

	UpdateUserIn struct {
		Name     string              `json:"name"`
		Surname  string              `json:"surname"`
		Role     domain.UserRole     `json:"role"`
		Position domain.UserPosition `json:"position"`
	}
)
