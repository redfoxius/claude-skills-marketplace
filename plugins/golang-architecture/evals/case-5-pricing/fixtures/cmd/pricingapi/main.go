package main

import (
	"fmt"

	"pricingapi/internal/gateway"
	"pricingapi/internal/pricing"
)

func main() {
	provider := gateway.NewOpenExchangeRatesClient("test-key")
	svc := pricing.NewService(provider)

	converted, err := svc.Convert(100, "USD", "EUR")
	if err != nil {
		panic(err)
	}

	fmt.Println(converted)
}
