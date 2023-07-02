package dao

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Topic struct {
	bun.BaseModel `bun:"table:topics"`
	ID            uuid.UUID  `bun:"type:uuid,default:uuid_generate_v4()"`
	Media         string     `bun:",notnull"`
	Title         string     `bun:",unique,notnull"`
	NSFW          bool       `bun:",notnull,default:false"`
	MaxThreads    int64      `bun:",notnull,default:64"`
	CreatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `bun:",nullzero,soft_delete"`
}
