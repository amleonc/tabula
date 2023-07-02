package main

import "github.com/amleonc/tabula/internal/web"

func main() {
	server := web.New()
	server.Start()
}
