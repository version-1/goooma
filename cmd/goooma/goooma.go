package main

import (
	"context"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/version-1/goooma"
)

func main() {
	g, err := goooma.New()
	if err != nil {
		log.Fatal(err)
	}

	g.Run(context.Background())
}
