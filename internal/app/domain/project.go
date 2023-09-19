package domain

import "time"

type (
	Project struct {
		ID          int32      `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		CoverId     string     `json:"coverId,omitempty"`
		Link        string     `json:"link,omitempty"`
		IsActive    bool       `json:"isActive"`
		TeamID      *int32     `json:"teamID,omitempty"`
		StartedAt   time.Time  `json:"startedAt"`
		EndedAt     *time.Time `json:"endedAt,omitempty"`
	}

	ProjectParticipant struct {
		UserID    int32        `json:"userID"`
		ProjectID int32        `json:"projectID"`
		Name      string       `json:"name"`
		Surname   string       `json:"surname"`
		Username  string       `json:"username"`
		Role      UserRole     `json:"role"`
		Position  UserPosition `json:"position"`
	}
)
