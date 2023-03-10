package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRole int16

const (
	_ UserRole = iota
	UserRoleGlobalAdmin
	UserRoleAdmin
	UserRoleModerator
	UserRoleUser
)

type (
	User struct {
		Id        int64      `json:"id"`
		Username  string     `json:"username"`
		Email     string     `json:"email"`
		Password  string     `json:"-"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		Role      UserRole   `json:"role"`
	}

	CreateUserRequest struct {
		Id       int64    `json:"id"`
		Username string   `json:"username"`
		Email    string   `json:"email"`
		Password string   `json:"-"`
		Role     UserRole `json:"role"`
	}
	CreateUserResponse struct {
		User *User `json:"data"`
	}

	GetUserRequest struct {
		UserId int64 `json:"-"`
	}
	GetUserResponse struct {
		User *User `json:"data"`
	}

	UpdateUserRequest struct {
		UserId   int64     `json:"userId"`
		Username *string   `json:"username"`
		Email    *string   `json:"email"`
		Password *string   `json:"password"`
		Role     *UserRole `json:"role"`
	}
	UpdateUserResponse struct {
		User *User `json:"data"`
	}

	DeleteUserRequest struct {
		UserId int64 `json:"userId"`
	}
	DeleteUserResponse struct{}
)

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
