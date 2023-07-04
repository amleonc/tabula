package internal

import (
	"context"

	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/db"
)

type Repo struct{}

func (Repo) Insert(ctx context.Context, m *dao.Media) (*dao.Media, error) {
	err := db.Insert(ctx, m)
	return m, err
}

func (Repo) Select(ctx context.Context, m *dao.Media) (*dao.Media, error) {
	r, err := db.SelectOne(ctx, m, "id", m.ID, "*")
	if err != nil {
		return nil, err
	}
	return r, nil
}
