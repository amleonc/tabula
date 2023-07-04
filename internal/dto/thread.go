package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type Thread struct {
	ID        uuid.UUID
	Media     *Media
	Topic     uuid.UUID
	User      uuid.UUID
	Title     string
	Body      string
	NSFW      bool
	Comments  []*Comment
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
