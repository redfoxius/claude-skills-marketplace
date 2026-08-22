package main

import (
	"log"

	"inventoryapi/internal/inventory"
)

func main() {
	svc := inventory.NewService()

	if err := svc.Reserve("sku-001", 3); err != nil {
		log.Fatal(err)
	}
}
