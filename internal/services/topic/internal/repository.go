package internal

import (
	"context"

	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/db"
	"github.com/gofrs/uuid"
)

// Types

type Repository struct{}

// Methods

func (Repository) Create(ctx context.Context, t *dao.Topic) error {
	return db.Insert(ctx, t)
}

func (Repository) SelectOneByID(ctx context.Context, id uuid.UUID) (*dao.Topic, error) {
	return db.SelectOne(ctx, &dao.Topic{}, "id", id, "*")
}
