package dao

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Thread struct {
	bun.BaseModel `bun:"table:threads"`
	ID            uuid.UUID  `bun:"type:uuid,default:uuid_generate_v4()"`
	Media         string     `bun:",notnull"`
	Topic         uuid.UUID  `bun:",notnull"`
	User          uuid.UUID  `bun:",notnull"`
	Title         string     `bun:",unique,notnull"`
	Body          string     `bun:",notnull"`
	NSFW          bool       `bun:",notnull,default:false"`
	CreatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `bun:",nullzero,soft_delete"`
}
