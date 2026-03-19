package database

import (
	"context"
	"log"
	"os"

	"authflow/ent"

	_ "github.com/lib/pq"
)

func NewClient() *ent.Client {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}

	// Auto-migrate (creates tables)
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema: %v", err)
	}

	return client
}
