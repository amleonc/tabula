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
	lookupCol string,
	searchVal any,
	selectCols ...string,
) (T, error) {
	err := db.NewSelect().
		Model(target).
		Column(selectCols...).
		Where("? = ?", bun.Ident(lookupCol), searchVal).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func SelectMultiple[T Record](
	ctx context.Context,
	target []T,
	lookupCol string,
	searchVal any,
	limit int,
	selectCols ...string,
) error {
	return db.NewSelect().
		Model(target).
		Column(selectCols...).
		Where("? > ?", bun.Ident(lookupCol), searchVal).
		Limit(limit).
		Scan(ctx)
}
