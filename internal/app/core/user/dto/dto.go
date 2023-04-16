package dto

import (
	"web-studio-backend/internal/app/core/user"	
)

type(
	UserCreate struct {
		Name     string            `json:"name"`
		Surname  string            `json:"surname"`
		Login    string            `json:"login"`
		Password string            `json:"password"`
		Role     user.UserRole     `json:"role"`
		Position user.UserPosition `json:"position"`
	}

	UserGet struct {
		UserId int16 `json:"-"`
	}

	UserUpdate struct {
		UserId   int16             `json:"-"`
		Name     string            `json:"name"`
		Surname  string            `json:"surname"`
		Role     user.UserRole     `json:"role"`
		Position user.UserPosition `json:"position"`
	}

	UserDelete struct {
		UserId int16 `json:"-"`
	}

	UserObject struct {
		User *user.User `json:"data"`
	}
)