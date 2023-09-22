package domain

import (
	"fmt"
	"net/url"
	"time"

	"web-studio-backend/internal/app/domain/apperr"
)

type (
	Project struct {
		ID           int32     `json:"id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		IsActive     bool      `json:"isActive"`
		Technologies []string  `json:"technologies,omitempty"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`

		CoverId *string    `json:"-"`
		Link    *string    `json:"link,omitempty"`
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
	var validations []apperr.ValidationError

	if p.Title == "" {
		validations = append(validations, apperr.ValidationError{
			Message: "Title cannot be empty.",
			Field:   "title",
		})
	}
	if len(p.Title) > 128 {
		validations = append(validations, apperr.ValidationError{
			Message: fmt.Sprintf("Title must be less than %d characters.", 128),
			Field:   "title",
		})
	}

	if len(p.Description) > 10000 {
		validations = append(validations, apperr.ValidationError{
			Message: fmt.Sprintf("Description must be less than %d characters.", 10000),
			Field:   "description",
		})
	}

	if p.Link != nil {
		_, err := url.ParseRequestURI(*p.Link)
		if err != nil {
			validations = append(validations, apperr.ValidationError{
				Message: "Link has invalid format.",
				Field:   "link",
			})
		}
	}

	if len(validations) > 0 {
		return apperr.NewValidationError(validations, "")
	}

	return nil
}

func (pp *ProjectParticipant) Validate() error {
	if pp.Role.String() == "" {
		return apperr.NewInvalidRequest("Unknown participant role.", "role")
	}

	if pp.Position.String() == "" {
		return apperr.NewInvalidRequest("Unknown participant position.", "position")
	}

	return nil
}
