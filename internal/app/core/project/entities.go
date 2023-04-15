package project

import (
	"time"
)

type Project struct {
	Id          int16     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverId     *string    `json:"coverId,omitempty"`
	StartedAt   time.Time `json:"startedAt"`
	EndedAt     time.Time `json:"endedAt,omitempty"`
	Link        *string   `json:"link,omitempty"`
}