package dto

import "web-studio-backend/internal/app/domain"

type (
	CreateTeamRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	UpdateTeamRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
)

func (r *CreateTeamRequest) ToDomain() *domain.Team {
	if r == nil {
		return nil
	}

	return &domain.Team{
		Title:       r.Title,
		Description: r.Description,
	}
}

func (r *UpdateTeamRequest) ToDomain(teamID int32) *domain.Team {
	if r == nil {
		return nil
	}

	return &domain.Team{
		ID:          teamID,
		Title:       r.Title,
		Description: r.Description,
	}
}
