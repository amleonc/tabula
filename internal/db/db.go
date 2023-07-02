package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/amleonc/tabula/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var db = initDB()

func initDB() *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DatabaseUser(),
		config.DatabasePassword(),
		config.DatabaseHost(),
		config.DatabasePort(),
		config.DatabaseName(),
		config.DatabaseSSLMode())
	log.Println(dsn)
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	if config.AppEnv() == "dev" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	return db
}
