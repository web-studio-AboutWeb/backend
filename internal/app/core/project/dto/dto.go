package dto

import (
	"time"
	"web-studio-backend/internal/app/core/user"	
	"web-studio-backend/internal/app/core/project"	
)

type (
	ProjectGet struct {
		ProjectId int16 `json:"-"`
	}

	ProjectParticipantsGet struct {
		ProjectId int16  `json:"-"`
		Role      string `json:"-"`
	}
	
	ProjectParticipants struct {
		Participants []user.User `json:"data"`
	}

	ProjectCreate struct {
		Id          int64      `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		StartedAt   time.Time  `json:"startedAt"`
		EndedAt     time.Time  `json:"endedAt"`
		Link        *string    `json:"link,omitempty"`
	}
	
	ProjectObject struct {
		Project *project.Project `json:"data"`
	}

	ProjectUpdate struct {
		ProjectId   int16     `json:"-"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartedAt   time.Time `json:"startedAt"`
		EndedAt     time.Time `json:"endedAt"`
		Link        *string   `json:"link,omitempty"`
	}

	ProjectCover struct {
		ProjectId      int16  `json:"-"`
		CoverImageData []byte `json:"-"`
	}

	ProjectDelete struct {
		ProjectId int16 `json:"-"`
	}
)
