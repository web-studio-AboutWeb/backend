package domain

import (
	"fmt"
	"time"

	"web-studio-backend/internal/app/domain/apperror"
)

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
		ID          int32      `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		HasImage    bool       `json:"hasImage"`
		CreatedAt   time.Time  `json:"createdAt"`
		UpdatedAt   time.Time  `json:"updatedAt"`
		DisabledAt  *time.Time `json:"disabledAt,omitempty"`

		ImageID      string `json:"-"`
		ImageContent []byte `json:"-"`
	}

	TeamMember struct {
		UserID    int32        `json:"userID"`
		TeamID    int32        `json:"teamID"`
		Role      UserRole     `json:"role"`
		Position  UserPosition `json:"position"`
		CreatedAt time.Time    `json:"createdAt"`
	}
)

func (t *Team) Validate() error {
	if t.Title == "" {
		return apperror.NewInvalidRequest("Title cannot be empty.", "title")
	}

	if len(t.Description) > 512 {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Description length must be less than %d characters.", 512),
			"description",
		)
	}

	return nil
}
