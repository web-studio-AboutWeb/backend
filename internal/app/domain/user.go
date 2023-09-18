package domain

import (
	"crypto/sha512"
	"fmt"
	"time"

	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/pkg/strhelp"
)

type (
	UserRole     int16
	UserPosition int16
)

const (
	UserRoleUser UserRole = iota + 1
	UserRoleModerator
	UserRoleAdmin
	UserRoleGlobalAdmin
)

const (
	UserPositionFrontend UserPosition = iota + 1
	UserPositionBackend
	UserPositionTeamLead
	UserPositionManager
	UserPositionMarketer
	UserPositionDevOps
)

type User struct {
	ID              int16        `json:"id"`
	Name            string       `json:"name"`
	Surname         string       `json:"surname"`
	Login           string       `json:"login"`
	EncodedPassword string       `json:"-"`
	Salt            string       `json:"-"`
	CreatedAt       time.Time    `json:"createdAt"`
	Role            UserRole     `json:"role"`
	RoleName        string       `json:"roleName"`
	Position        UserPosition `json:"position"`
	PositionName    string       `json:"positionName"`
}

func (u *User) Validate() error {
	if u.Name == "" || len(u.Name) > 30 {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Name cannot be empty and must not exceed %d characters.", 30),
			"name",
		)
	}

	if u.Surname == "" || len(u.Surname) > 50 {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Surname cannot be empty and must not exceed %d characters.", 50),
			"surname",
		)
	}

	if u.Login == "" || len(u.Login) > 20 {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Login cannot be empty and must not exceed %d characters.", 20),
			"login",
		)
	}

	if u.EncodedPassword == "" || len(u.EncodedPassword) > 20 {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Password cannot be empty and must not exceed %d characters.", 20),
			"login",
		)
	}

	if u.Role.String() == "" {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Unknown role %d.", u.Role),
			"role",
		)
	}

	if u.Position.String() == "" {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Unknown position %d.", u.Position),
			"position",
		)
	}

	return nil
}

func (u *User) EncodePassword() error {
	if u.EncodedPassword == "" {
		return fmt.Errorf("empty password")
	}

	salt, err := strhelp.GenerateRandomString(32)
	if err != nil {
		return fmt.Errorf("generating random string: %w", err)
	}
	u.Salt = salt

	u.EncodedPassword = fmt.Sprintf("%x", sha512.Sum512([]byte(u.EncodedPassword+u.Salt)))

	return nil
}

func (u *User) ComparePassword(password string) bool {
	if password == "" || u.Salt == "" || u.EncodedPassword == "" {
		return false
	}

	passwordHashBytes := sha512.Sum512(append([]byte(password), u.Salt...))
	passwordHash := fmt.Sprintf("%x", passwordHashBytes)

	return passwordHash == u.EncodedPassword
}

func (ur UserRole) String() string {
	switch ur {
	case UserRoleUser:
		return "User"
	case UserRoleModerator:
		return "Moderator"
	case UserRoleAdmin:
		return "Admin"
	case UserRoleGlobalAdmin:
		return "Global admin"
	default:
		return ""
	}
}

func (up UserPosition) String() string {
	switch up {
	case UserPositionFrontend:
		return "Frontend"
	case UserPositionBackend:
		return "Backend"
	case UserPositionTeamLead:
		return "Team lead"
	case UserPositionManager:
		return "Manager"
	case UserPositionMarketer:
		return "Marketer"
	case UserPositionDevOps:
		return "DevOps"
	default:
		return ""
	}
}
