package internal

import (
	"context"

	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/db"
	"github.com/gofrs/uuid"
)

type Repository struct{}

func (Repository) Create(ctx context.Context, c *dao.Comment) error {

	q := db.NewInsertQuery().
		Model(c).
		Value(
			"is_op",
			`(CASE WHEN (SELECT "user" FROM threads WHERE id = ?) = ? THEN TRUE ELSE FALSE END)`,
			c.Thread, c.User).
		Returning("*")

	err := q.Scan(ctx)
	if err != nil {
		return err
	}

	return nil

}

func (Repository) SelectByThreadID(ctx context.Context, id uuid.UUID, limit int) ([]*dao.Comment, error) {
	var target []*dao.Comment

	err := db.SelectMultiple(ctx, &target, map[string]any{
		"thread": id.String(),
	}, limit, "*")

	return target, err
}
