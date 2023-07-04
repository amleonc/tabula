package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type Topic struct {
	ID         uuid.UUID  `json:"id,omitempty"`
	Media      *Media     `json:"media,omitempty"`
	User       uuid.UUID  `json:"user,omitempty"`
	Title      string     `json:"title,omitempty"`
	ShortTitle string     `json:"short_title,omitempty"`
	NSFW       bool       `json:"nsfw,omitempty"`
	MaxThreads int64      `json:"max_threads,omitempty"`
	Threads    []*Thread  `json:"threads,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
