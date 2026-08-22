package main

import (
	"context"
	"database/sql"
	"log"

	"subscriptionapi/internal/store"
	"subscriptionapi/internal/subscription"
)

func main() {
	db, err := sql.Open("postgres", "postgres://localhost/billing")
	if err != nil {
		log.Fatal(err)
	}

	st := store.NewPostgresStore(db)
	svc := subscription.NewService(st)

	if err := svc.Activate(context.Background(), "user_1", "pro"); err != nil {
		log.Fatal(err)
	}
}
