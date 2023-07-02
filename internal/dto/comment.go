package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	ID        uuid.UUID  `json:"id,omitempty"`
	Media     *Media     `json:"media,omitempty"`
	Thread    uuid.UUID  `json:"thread,omitempty"`
	User      uuid.UUID  `json:"user,omitempty"`
	Grip      string     `json:"grip,omitempty"`
	Body      string     `json:"body,omitempty"`
	Color     uint8      `json:"color,omitempty"`
	IsOP      bool       `json:"is_op,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
