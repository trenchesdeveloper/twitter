package main

import (
	"context"
	"log"

	"github.com/trenchesdeveloper/tweeter/config"
	"github.com/trenchesdeveloper/tweeter/postgres"
)

func main() {
	ctx := context.Background()

	// create config
	config := config.New()

	// create db
	db := postgres.New(ctx, config)

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("migrations ran successfully")
}
