package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Password  string     `json:"password,omitempty"`
	Role      int64      `json:"role,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
