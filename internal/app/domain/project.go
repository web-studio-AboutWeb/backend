package domain

import (
	"fmt"
	"net/url"
	"time"

	"web-studio-backend/internal/app/domain/apperror"
)

type (
	Project struct {
		ID          int32     `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		CoverId     string    `json:"coverId,omitempty"`
		Link        string    `json:"link,omitempty"`
		IsActive    bool      `json:"isActive"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt,omitempty"`

		TeamID  *int32     `json:"teamID,omitempty"`
		EndedAt *time.Time `json:"endedAt,omitempty"`
	}

	ProjectParticipant struct {
		UserID    int32        `json:"userID"`
		ProjectID int32        `json:"projectID"`
		Name      string       `json:"name"`
		Surname   string       `json:"surname"`
		Username  string       `json:"username"`
		Role      UserRole     `json:"role"`
		Position  UserPosition `json:"position"`
		CreatedAt time.Time    `json:"createdAt"`
		UpdatedAt time.Time    `json:"updatedAt"`
	}
)

func (p *Project) Validate() error {
	if p.Title == "" {
		return apperror.NewInvalidRequest("Title cannot be empty.", "title")
	}
	if len(p.Title) > 128 {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Title must be less than %d characters.", 128),
			"title",
		)
	}

	if len(p.Description) > 2048 {
		return apperror.NewInvalidRequest(
			fmt.Sprintf("Description must be less than %d characters.", 2048),
			"description",
		)
	}

	if p.Link != "" {
		_, err := url.ParseRequestURI(p.Link)
		if err != nil {
			return apperror.NewInvalidRequest("Link has invalid format.", "link")
		}
	}

	return nil
}

func (pp *ProjectParticipant) Validate() error {
	if pp.Role.String() == "" {
		return apperror.NewInvalidRequest("Unknown participant role.", "role")
	}

	if pp.Position.String() == "" {
		return apperror.NewInvalidRequest("Unknown participant position.", "position")
	}

	return nil
}
