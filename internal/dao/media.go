package dao

import (
	"time"

	"github.com/uptrace/bun"
)

type Media struct {
	bun.BaseModel `bun:"table:media"`
	ID            string     `bun:",pk"`
	Type          string     `bun:",notnull"`
	Extension     string     `bun:",notnull"`
	Blacklist     bool       `bun:",nullzero,notnull,default:false"` // in case someone wants to upload a forbidden image.
	CreatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `bun:",nullzero,soft_delete"`
}
