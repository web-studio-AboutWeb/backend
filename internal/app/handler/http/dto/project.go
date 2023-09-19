package dto

import (
	"web-studio-backend/internal/app/domain"
)

type (
	CreateProjectRequest struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Link        *string `json:"link,omitempty"`
		TeamID      *int32  `json:"teamID,omitempty"`
	}

	UpdateProjectRequest struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Link        *string `json:"link,omitempty"`
		TeamID      *int32  `json:"teamID,omitempty"`
	}

	AddProjectParticipantRequest struct {
		UserID   int32               `json:"userID"`
		Role     domain.UserRole     `json:"role"`
		Position domain.UserPosition `json:"position"`
	}

	UpdateProjectParticipantRequest struct {
		Role     domain.UserRole     `json:"role"`
		Position domain.UserPosition `json:"position"`
	}
)

func (r *CreateProjectRequest) ToDomain() *domain.Project {
	if r == nil {
		return nil
	}

	return &domain.Project{
		Title:       r.Title,
		Description: r.Description,
		Link:        r.Link,
		TeamID:      r.TeamID,
	}
}

func (r *UpdateProjectRequest) ToDomain(projectID int32) *domain.Project {
	if r == nil {
		return nil
	}

	return &domain.Project{
		ID:          projectID,
		Title:       r.Title,
		Description: r.Description,
		Link:        r.Link,
		TeamID:      r.TeamID,
	}
}

func (r *AddProjectParticipantRequest) ToDomain(projectID int32) *domain.ProjectParticipant {
	if r == nil {
		return nil
	}

	return &domain.ProjectParticipant{
		UserID:    r.UserID,
		ProjectID: projectID,
		Role:      r.Role,
		Position:  r.Position,
	}
}

func (r *UpdateProjectParticipantRequest) ToDomain(projectID, userID int32) *domain.ProjectParticipant {
	if r == nil {
		return nil
	}

	return &domain.ProjectParticipant{
		UserID:    userID,
		ProjectID: projectID,
		Role:      r.Role,
		Position:  r.Position,
	}
}
