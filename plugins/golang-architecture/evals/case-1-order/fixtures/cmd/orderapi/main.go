package main

import (
	"log"

	"orderapi/internal/order"
	"orderapi/internal/postgres"
)

func main() {
	repo := postgres.NewPGOrderRepository("postgres://localhost/orders")

	svc, err := order.NewService(repo, "postgres://localhost/orders")
	if err != nil {
		log.Fatal(err)
	}

	if err := svc.MarkShipped("ord_123"); err != nil {
		log.Fatal(err)
	}
}
