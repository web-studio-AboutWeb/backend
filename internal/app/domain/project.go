package domain

import "time"

type (
	Project struct {
		ID          int32      `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		CoverId     string     `json:"coverId,omitempty"`
		Link        string     `json:"link,omitempty"`
		TeamID      *int32     `json:"teamID,omitempty"`
		StartedAt   time.Time  `json:"startedAt"`
		EndedAt     *time.Time `json:"endedAt,omitempty"`
	}

	ProjectParticipant struct {
		ProjectID int32        `json:"projectID"`
		UserID    int32        `json:"userID"`
		Role      UserRole     `json:"role"`
		Position  UserPosition `json:"position"`
	}
)
