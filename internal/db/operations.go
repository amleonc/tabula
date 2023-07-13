package db

import (
	"context"
	"log"

	"github.com/amleonc/tabula/internal/dao"
	"github.com/uptrace/bun"
)

type Record interface {
	*dao.Comment |
		*dao.Media |
		*dao.Thread |
		*dao.Topic |
		*dao.User
}

func CreateTable[T Record](ctx context.Context, table T) error {
	if table == nil {
		log.Println("skipping model")
	}
	_, err := db.NewCreateTable().
		Model(table).
		IfNotExists().
		Exec(ctx)
	return err
}

func Insert[T Record](ctx context.Context, thing T) error {
	return db.NewInsert().
		Model(thing).
		Scan(ctx, thing)
}

func SelectOne[T Record](
	ctx context.Context,
	target T,
	filter Filter,
	selectCols ...string,
) (T, error) {
	query := db.NewSelect().
		Model(target).
		Column(selectCols...)
	for k, v := range filter {
		query.Where("? = ?", bun.Ident(k), v)
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return target, nil
}

const (
	defaultLimit = 40
)

func SelectMultiple[T Record](
	ctx context.Context,
	target *[]T,
	filter map[string]any,
	limit int,
	selectCols ...string,
) error {
	query := db.NewSelect().
		Model(target).
		Column(selectCols...)

	for k, v := range filter {
		query.Where("? = ?", bun.Ident(k), v)
	}

	if limit == 0 {
		limit = defaultLimit
	}
	query.Limit(limit)

	return query.Scan(ctx)
}

func NewSelectQuery() *bun.SelectQuery {
	return db.NewSelect()
}

func NewInsertQuery() *bun.InsertQuery {
	return db.NewInsert()
}
