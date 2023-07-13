package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type Thread struct {
	ID        uuid.UUID  `json:"id,omitempty"`
	Media     *Media     `json:"media,omitempty"`
	Topic     uuid.UUID  `json:"topic,omitempty"`
	User      uuid.UUID  `json:"-"`
	Title     string     `json:"title,omitempty"`
	Body      string     `json:"body,omitempty"`
	NSFW      bool       `json:"nsfw,omitempty"`
	Comments  []*Comment `json:"comments,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
