package main

import (
	"context"
	"log"

	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/db"
)

func main() {
	ctx := context.Background()
	var err error
	err = db.CreateTable(ctx, &dao.Comment{})
	handle(err)
	err = db.CreateTable(ctx, &dao.Media{})
	handle(err)
	err = db.CreateTable(ctx, &dao.Thread{})
	handle(err)
	err = db.CreateTable(ctx, &dao.Topic{})
	handle(err)
	err = db.CreateTable(ctx, &dao.User{})
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
