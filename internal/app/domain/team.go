package domain

import "time"

type UserPosition int16

const (
	UserPositionFrontend UserPosition = iota + 1
	UserPositionBackend
	UserPositionManager
	UserPositionMarketer
	UserPositionDevOps
)

func (up UserPosition) String() string {
	switch up {
	case UserPositionFrontend:
		return "Frontend"
	case UserPositionBackend:
		return "Backend"
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

type (
	Team struct {
		ID         int32      `json:"id"`
		Title      string     `json:"title"`
		CreatedAt  time.Time  `json:"createdAt"`
		UpdatedAt  time.Time  `json:"updatedAt"`
		DisabledAt *time.Time `json:"disabledAt,omitempty"`
	}

	TeamMember struct {
		UserID   int32        `json:"userID"`
		TeamID   int32        `json:"teamID"`
		Role     UserRole     `json:"role"`
		Position UserPosition `json:"position"`
	}
)
