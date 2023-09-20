package domain

import "time"

type Document struct {
	ID               int32     `json:"id"`
	OriginalFilename string    `json:"originalFilename"`
	FileID           string    `json:"-"`
	UserID           int32     `json:"userID"`
	CreatedAt        time.Time `json:"createdAt"`
	MimeType         string    `json:"mimeType"`
	SizeBytes        int32     `json:"sizeBytes"` // Size in bytes
	Content          []byte    `json:"-"`
}
