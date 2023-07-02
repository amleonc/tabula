package dto

import (
	"io"
	"time"
)

type Media struct {
	ID                string            `json:"id,omitempty"`
	Filename          string            `json:"url,omitempty"`
	ThumbnailFilename string            `json:"thumbnail_filename,omitempty"`
	Type              string            `json:"type,omitempty"`
	Format            string            `json:"format,omitempty"`
	Bytes             io.ReadSeekCloser `json:"file,omitempty"`
	Blacklist         bool              `json:"blacklist,omitempty"`
	CreatedAt         *time.Time        `json:"created_at,omitempty"`
	UpdatedAt         *time.Time        `json:"updated_at,omitempty"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty"`
}
