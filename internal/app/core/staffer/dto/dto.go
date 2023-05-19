package dto

import (
	"web-studio-backend/internal/app/core/staffer"
)

type (
	StafferCreate struct {
		UserId    int16                   `json:"user_id"`
		ProjectId int16                   `json:"project_id"`
		Position  staffer.StafferPosition `json:"position"`
	}

	StafferGet struct {
		StafferId int16 `json:"-"`
	}
	StaffersGet struct {
		ProjectId int16 `json:"-"`
	}

	StafferUpdate struct {
		StafferId int16                   `json:"-"`
		ProjectId int16                   `json:"project_id"`
		Position  staffer.StafferPosition `json:"position"`
	}

	StafferDelete struct {
		StafferId int16 `json:"-"`
	}

	StafferObject struct {
		Staffer *staffer.Staffer `json:"data"`
	}
	StaffersObject struct {
		Staffers []staffer.Staffer `json:"data"`
	}
)
