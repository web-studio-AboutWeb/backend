package dto

import (
	"time"

	"web-studio-backend/internal/app/core/project"
	"web-studio-backend/internal/app/core/staffer"
	"web-studio-backend/internal/app/core/user"
)

type (
	ProjectGet struct {
		ProjectId int16 `json:"-"`
	}

	ProjectStaffersGet struct {
		ProjectId int16  `json:"-"`
		Role      string `json:"-"`
	}

	ProjectStaffers struct {
		Staffers []ProjectStaffer `json:"data"`
	}

	ProjectStaffer struct {
		Id        int16                   `json:"id"`
		User      user.User               `json:"user"`
		ProjectId int16                   `json:"project_id"`
		Position  staffer.StafferPosition `json:"position"`
	}

	ProjectCreate struct {
		Id          int64     `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartedAt   time.Time `json:"startedAt"`
		EndedAt     time.Time `json:"endedAt"`
		Link        *string   `json:"link,omitempty"`
	}

	ProjectObject struct {
		Project *project.Project `json:"data"`
	}

	ProjectsObject struct {
		Projects []project.Project `json:"data"`
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
