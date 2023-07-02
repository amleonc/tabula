package dao

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Comment struct {
	bun.BaseModel `bun:"table:comments"`
	ID            uuid.UUID  `bun:"type:uuid,default:uuid_generate_v4()"`
	Media         string     `bun:",nullzero"`
	Thread        uuid.UUID  `bun:",notnull"`
	User          uuid.UUID  `bun:",notnull"`
	Grip          string     `bun:",notnull"`
	Body          string     `bun:",unique,notnull"`
	Color         uint8      `bun:",notnull"`
	IsOP          bool       `bun:",nullzero,notnull,default:false"`
	CreatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `bun:",nullzero,soft_delete"`
}
