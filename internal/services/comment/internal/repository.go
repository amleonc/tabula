package internal

import (
	"context"

	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/db"
	"github.com/gofrs/uuid"
)

type Repository struct{}

func (Repository) Create(ctx context.Context, c *dao.Comment) error {
	return db.Insert(ctx, c)
}

func (Repository) SelectByThreadID(ctx context.Context, id uuid.UUID, limit int) ([]*dao.Comment, error) {
	var target []*dao.Comment
	err := db.SelectMultiple(ctx, target, "thread", id, limit, "*")
	return target, err
}
