package main

import (
	"log"

	"billingapi/internal/billing"
)

func main() {
	brokers := []string{"localhost:9092"}

	svc := billing.NewService(brokers)

	inv := billing.Invoice{
		ID:          "inv_1",
		CustomerID:  "cus_1",
		AmountCents: 4200,
	}

	if err := svc.IssueInvoice(inv); err != nil {
		log.Fatal(err)
	}
}
