package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	UserRole     int16
	UserPosition int16
)

const (
	_ UserRole = iota
	UserRoleUser
	UserRoleModerator
	UserRoleAdmin
	UserRoleGlobalAdmin
)

const (
	_ UserPosition = iota
	UserPositionFrontend
	UserPositionBackend
	UserPositionTeamLead
	UserPositionManager
	UserPositionMarketer
	UserPositionDevOps
)

type User struct {
	ID        int16        `json:"id"`
	Name      string       `json:"name"`
	Surname   string       `json:"surname"`
	Login     string       `json:"-"`
	Password  string       `json:"-"`
	CreatedAt time.Time    `json:"createdAt"`
	Role      UserRole     `json:"role"`
	Position  UserPosition `json:"position"`
}

func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
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
		return "None"
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
		return "None"
	}
}
