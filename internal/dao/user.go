package dao

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            uuid.UUID  `bun:"type:uuid,default:gen_random_uuid()"`
	Name          string     `bun:",unique,notnull"`
	Password      string     `bun:",nullzero,notnull"`
	Role          int64      `bun:",nullzero,notnull,default:3"`
	CreatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `bun:",nullzero,soft_delete"`
}
