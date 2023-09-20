package domain

import "time"

type Document struct {
	ID        int32     `json:"id"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"createdAt"`
	MimeType  string    `json:"mimeType"`
	Size      int32     `json:"size"` // Size in bytes
}
