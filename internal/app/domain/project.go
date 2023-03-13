package domain

import "time"

type (
	Project struct {
		Id          int16     `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		CoverId     string    `json:"coverId,omitempty"`
		StartedAt   time.Time `json:"startedAt"`
		EndedAt     time.Time `json:"endedAt,omitempty"`
		Link        string    `json:"link,omitempty"`
	}

	GetProjectRequest struct {
		ProjectId int16 `json:"-"`
	}
	GetProjectResponse struct {
		Project *Project `json:"data"`
	}
	GetProjectParticipantsRequest struct {
		ProjectId int16  `json:"-"`
		Role      string `json:"-"`
	}
	GetProjectParticipantsResponse struct {
		Participants []User `json:"data"`
	}

	CreateProjectRequest struct {
		Id          int64      `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		StartedAt   time.Time  `json:"startedAt"`
		EndedAt     *time.Time `json:"endedAt"`
		Link        *string    `json:"link,omitempty"`
	}
	CreateProjectResponse struct {
		Project *Project `json:"data"`
	}

	UpdateProjectRequest struct {
		ProjectId   int16     `json:"-"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartedAt   time.Time `json:"startedAt"`
		EndedAt     time.Time `json:"endedAt"`
		Link        string    `json:"link,omitempty"`
	}
	UpdateProjectResponse struct {
		Project *Project `json:"data"`
	}

	UploadProjectCoverRequest struct {
		ProjectId      int16  `json:"-"`
		CoverImageData []byte `json:"-"`
	}
	UploadProjectCoverResponse struct {
		Project *Project `json:"data"`
	}

	DeleteProjectRequest struct {
		ProjectId int16 `json:"-"`
	}
	DeleteProjectResponse struct{}
)
