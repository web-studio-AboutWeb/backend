package user

import (
	"time"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserRole     string
	UserPosition string
)

const (
	UserRoleSuperAdmin = "super admin"
	UserRoleAdmin       = "admin"
	UserRoleModerator   = "moderator"
	UserRoleUser        = "user"
)

type User struct {
	Id        int16        `json:"id"`
	Name      string       `json:"name"`
	Surname   string       `json:"surname"`
	Login     string       `json:"-"`
	Password  string       `json:"-"`
	CreatedAt time.Time    `json:"createdAt"`
	Role      UserRole     `json:"role"`
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
