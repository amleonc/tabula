package internal

import (
	"context"

	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/db"
	"github.com/gofrs/uuid"
)

type Repostiroy struct{}

func (Repostiroy) Create(ctx context.Context, t *dao.Thread) error {
	return db.Insert(ctx, t)
}

func (Repostiroy) SelectByID(ctx context.Context, id uuid.UUID) (*dao.Thread, error) {
	return db.SelectOne(ctx, &dao.Thread{}, "id", id, "*")
}

func (Repostiroy) SelectByTopicID(ctx context.Context, id uuid.UUID, limit int) ([]*dao.Thread, error) {
	t := make([]*dao.Thread, 0)
	err := db.SelectMultiple(
		ctx,
		t,
		"topic",
		id,
		limit,
		"*",
	)
	return t, err
}
