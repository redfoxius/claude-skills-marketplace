package main

import (
	"log"

	"shippingapi/internal/carrier"
	"shippingapi/internal/shipping"
)

func main() {
	fedex := carrier.NewFedExClient("test-key")
	svc := shipping.NewService(fedex)

	trackingID, err := svc.ShipOrder("ord_1", 1200, "10001")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(trackingID)
}
