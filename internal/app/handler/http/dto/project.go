package dto

import (
	"web-studio-backend/internal/app/domain"
)

type (
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
