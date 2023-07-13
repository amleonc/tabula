package internal

import (
	"context"

	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/db"
)

type Repository struct{}

func (Repository) Create(ctx context.Context, u *dao.User) error {
	err := db.Insert(ctx, u)
	return err
}

func (Repository) SelectByName(ctx context.Context, n string) (*dao.User, error) {
	u := new(dao.User)
	var err error
	u, err = db.SelectOne[*dao.User](
		ctx,
		u,
		db.Filter{"name": n},
		"*",
	)
	return u, err
}
