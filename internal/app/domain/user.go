package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	UserRole     string
	UserPosition string
)

const (
	UserRoleGlobalAdmin = "global admin"
	UserRoleAdmin       = "admin"
	UserRoleModerator   = "moderator"
	UserRoleUser        = "user"
)

const (
	UserPositionFrontendDev = "frontend"
	UserPositionBackendDev  = "backend"
	UserPositionTeamLead    = "teamlead"
	UserPositionManager     = "manager"
	UserPositionMarketer    = "marketer"
)

type (
	User struct {
		Id        int16        `json:"id"`
		Name      string       `json:"name"`
		Surname   string       `json:"surname"`
		Login     string       `json:"-"`
		Password  string       `json:"-"`
		CreatedAt time.Time    `json:"createdAt"`
		Role      UserRole     `json:"role"`
		Position  UserPosition `json:"position"`
	}

	CreateUserRequest struct {
		Id       int64        `json:"id"`
		Name     string       `json:"name"`
		Surname  string       `json:"surname"`
		Login    string       `json:"login"`
		Password string       `json:"password"`
		Role     UserRole     `json:"role"`
		Position UserPosition `json:"position"`
	}
	CreateUserResponse struct {
		User *User `json:"data"`
	}

	GetUserRequest struct {
		UserId int16 `json:"-"`
	}
	GetUserResponse struct {
		User *User `json:"data"`
	}

	UpdateUserRequest struct {
		UserId   int16        `json:"-"`
		Name     string       `json:"name"`
		Surname  string       `json:"surname"`
		Role     UserRole     `json:"role"`
		Position UserPosition `json:"position"`
	}
	UpdateUserResponse struct {
		User *User `json:"data"`
	}

	DeleteUserRequest struct {
		UserId int16 `json:"-"`
	}
	DeleteUserResponse struct{}
)

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
